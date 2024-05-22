package controller

import (
	"encoding/json"
	"maden/pkg/etcd"
	"maden/pkg/shared"
	"maden/pkg/scheduler"

	"log"

	"github.com/google/uuid"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

// Create
func handleDeploymentCreate(kv *mvccpb.KeyValue) {
	log.Printf("New deployment created: %s", string(kv.Value))

	var deployment shared.Deployment
	if err := json.Unmarshal(kv.Value, &deployment); err != nil {
		log.Printf("Failed to unmarshal deployment: %v", err)
		return
	}

	createAndSchedulePodsFromDeployment(&deployment, -1)
}

func createAndSchedulePodsFromDeployment(deployment *shared.Deployment, limit int) {
	maxPods := getMaxPods(deployment.Replicas, limit)

	for i := 0; i < maxPods; i++ {
		pod := createPodFromTemplate(deployment.Template, deployment.Name, deployment.ID)
		if err := etcd.CreatePod(pod); err != nil {
			log.Printf("Failed to create pod: %v", err)
			return
		}
		if err := scheduler.SchedulePod(pod); err != nil {
			log.Printf("Failed to schedule pod: %v", err)
			return
		}
	}
}

func createPodFromTemplate(template shared.PodTemplate, podName string, deploymentID string) *shared.Pod {
	podID := podName + "-" + uuid.New().String()
	pod := &shared.Pod{
		ID: podID,
		Name: podName,
		DeploymentID: deploymentID,
		Status: shared.PodPending,
		NodeID: "",
		Resources: template.Spec.Resources,
		Affinity: template.Spec.Affinity,
		AntiAffinity: template.Spec.AntiAffinity,
		Tolerations: template.Spec.Tolerations,
	}
	return pod
}

// Update
func handleDeploymentUpdate(oldKV *mvccpb.KeyValue, newKV *mvccpb.KeyValue) {
	log.Printf("Deployment updated: %s, %v", string(oldKV.Value), string(newKV.Value))
	
	var oldDeployment shared.Deployment
	if err := json.Unmarshal(oldKV.Value, &oldDeployment); err != nil {
		log.Printf("Failed to unmarshal old deployment: %v", err)
		return
	}

	var newDeployment shared.Deployment
	if err := json.Unmarshal(newKV.Value, &newDeployment); err != nil {
		log.Printf("Failed to unmarshal new deployment: %v", err)
		return
	}

	if oldDeployment.Replicas != newDeployment.Replicas {
		handleDeploymentReplicasUpdate(&oldDeployment, &newDeployment)
	}
	if !arePodTemplatesEqual(oldDeployment.Template, newDeployment.Template) {
		handleDeploymentTemplateUpdate(oldDeployment.ID, &newDeployment)
	}
}

func handleDeploymentReplicasUpdate(oldDeployment *shared.Deployment, newDeployment *shared.Deployment) {
	difference := newDeployment.Replicas - oldDeployment.Replicas
	if difference > 0 {
		createAndSchedulePodsFromDeployment(newDeployment, difference)
	} else if difference < 0 {
		deletePodsByDeploymentID(newDeployment.ID, -difference)
	}
}

func handleDeploymentTemplateUpdate(oldDeploymentID string, newDeployment *shared.Deployment) {
	deletePodsByDeploymentID(oldDeploymentID, -1)
	createAndSchedulePodsFromDeployment(newDeployment, -1)
}

// Delete
func handleDeploymentDelete(kv *mvccpb.KeyValue) {
	log.Printf("Deployment deleted: %s", string(kv.Value))

	var deployment shared.Deployment
	if err := json.Unmarshal(kv.Value, &deployment); err != nil {
		log.Printf("Failed to unmarshal deployment: %v", err)
		return
	}

	deletePodsByDeploymentID(deployment.ID, -1)
}

func deletePodsByDeploymentID(deploymentID string, limit int) {
	pods, err := etcd.GetPodsByDeploymentID(deploymentID)
	if err != nil {
		log.Printf("Failed to get pods by deployment ID: %v", err)
		return
	}

	maxPods := getMaxPods(len(pods), limit)
	for _, pod := range pods[:maxPods] {
		err := etcd.DeletePod(pod.ID)
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