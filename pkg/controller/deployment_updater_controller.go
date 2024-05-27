package controller

import (
	"maden/pkg/etcd"
	"maden/pkg/orchestrator"
	"maden/pkg/shared"

	"encoding/json"
	"log"

	"github.com/google/uuid"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type DefaultDeploymentUpdaterController struct {
	Repo etcd.PodRepository
	Orchestrator orchestrator.PodOrchestrator
}

func NewDefaultDeploymentUpdaterController(
	repo etcd.PodRepository,
	orchestrator orchestrator.PodOrchestrator,
) DeploymentUpdaterController {
	return &DefaultDeploymentUpdaterController{Repo: repo, Orchestrator: orchestrator}
}

// Create
func (c *DefaultDeploymentUpdaterController) HandleDeploymentCreate(kv *mvccpb.KeyValue) {
	log.Printf("New deployment created: %s", string(kv.Value))

	var deployment shared.Deployment
	if err := json.Unmarshal(kv.Value, &deployment); err != nil {
		log.Printf("Failed to unmarshal deployment: %v", err)
		return
	}

	c.createAndSchedulePodsFromDeployment(&deployment, -1)
}

func (c *DefaultDeploymentUpdaterController) createAndSchedulePodsFromDeployment(deployment *shared.Deployment, limit int) {
	maxPods := getMaxPods(deployment.Replicas, limit)

	for i := 0; i < maxPods; i++ {
		pod := getPodFromTemplate(deployment.Template, deployment.Name, deployment.ID)
		if err := c.Orchestrator.OrchestratePodCreation(pod); err != nil {
			log.Printf("Failed to create pod: %v", err)
			return
		}
	}
}

func getPodFromTemplate(template shared.PodTemplate, podName string, deploymentID string) *shared.Pod {
	podID := podName + "-" + uuid.New().String()
	pod := &shared.Pod{
		ID: podID,
		Name: podName,
		DeploymentID: deploymentID,
		Status: shared.PodPending,
		NodeID: "",
		Containers: template.Spec.Containers,
		Resources: template.Spec.Resources,
		Affinity: template.Spec.Affinity,
		AntiAffinity: template.Spec.AntiAffinity,
		Tolerations: template.Spec.Tolerations,
		RestartPolicy: template.Spec.RestartPolicy,
	}
	return pod
}

// Update
func (c *DefaultDeploymentUpdaterController) HandleDeploymentUpdate(oldKv *mvccpb.KeyValue, newKv *mvccpb.KeyValue) {
	log.Printf("Deployment updated: %s, %v", string(oldKv.Value), string(newKv.Value))
	
	var oldDeployment shared.Deployment
	if err := json.Unmarshal(oldKv.Value, &oldDeployment); err != nil {
		log.Printf("Failed to unmarshal old deployment: %v", err)
		return
	}

	var newDeployment shared.Deployment
	if err := json.Unmarshal(newKv.Value, &newDeployment); err != nil {
		log.Printf("Failed to unmarshal new deployment: %v", err)
		return
	}

	if oldDeployment.Replicas != newDeployment.Replicas {
		c.handleDeploymentReplicasUpdate(&oldDeployment, &newDeployment)
	}
	if !arePodTemplatesEqual(oldDeployment.Template, newDeployment.Template) {
		c.handleDeploymentTemplateUpdate(oldDeployment.ID, &newDeployment)
	}
}

func (c *DefaultDeploymentUpdaterController) handleDeploymentReplicasUpdate(oldDeployment *shared.Deployment, newDeployment *shared.Deployment) {
	difference := newDeployment.Replicas - oldDeployment.Replicas
	if difference > 0 {
		c.createAndSchedulePodsFromDeployment(newDeployment, difference)
	} else if difference < 0 {
		c.deletePodsByDeploymentID(newDeployment.ID, -difference)
	}
}

func (c *DefaultDeploymentUpdaterController) handleDeploymentTemplateUpdate(oldDeploymentID string, newDeployment *shared.Deployment) {
	c.deletePodsByDeploymentID(oldDeploymentID, -1)
	c.createAndSchedulePodsFromDeployment(newDeployment, -1)
}

// Delete
func (c *DefaultDeploymentUpdaterController) HandleDeploymentDelete(kv *mvccpb.KeyValue) {
	log.Printf("Deployment deleted: %s", string(kv.Value))

	var deployment shared.Deployment
	if err := json.Unmarshal(kv.Value, &deployment); err != nil {
		log.Printf("Failed to unmarshal deployment: %v", err)
		return
	}

	c.deletePodsByDeploymentID(deployment.ID, -1)
}

func (c *DefaultDeploymentUpdaterController) deletePodsByDeploymentID(deploymentID string, limit int) {
	pods, err := c.Repo.GetPodsByDeploymentID(deploymentID)
	if err != nil {
		log.Printf("Failed to get pods by deployment ID: %v", err)
		return
	}

	maxPods := getMaxPods(len(pods), limit)
	for _, pod := range pods[:maxPods] {
		err := c.Orchestrator.OrchestratePodDeletion(&pod)
		if err != nil {
			log.Printf("Failed to delete pod %s: %v", pod.ID, err)
			return
		}
	}
}

func getMaxPods(podsCount int, limit int) int {
	var maxPods int
	if limit == -1 {
		maxPods = podsCount
	} else {
		maxPods = min(podsCount, limit)
	}
	return maxPods
}

func (c *DefaultDeploymentUpdaterController) HandleDeploymentRolloutRestart(deployment *shared.Deployment) error {
	pods, err := c.Repo.GetPodsByDeploymentID(deployment.ID)
	if err != nil {
		return err
	}

	for _, pod := range pods {
		err := c.Orchestrator.OrchestratePodDeletion(&pod)
		if err != nil {
			return err
		}

		newPod := getPodFromTemplate(deployment.Template, deployment.Name, deployment.ID)
		err = c.Orchestrator.OrchestratePodCreation(newPod)
		if err != nil {
			return err
		}
	}

	return nil
}