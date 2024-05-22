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

	createAndSchedulePodsFromDeployment(&deployment)
}

func createAndSchedulePodsFromDeployment(deployment *shared.Deployment) {
	for i := 0; i < deployment.Replicas; i++ {
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
func handleDeploymentUpdate(kv *mvccpb.KeyValue) {
	log.Printf("Deployment updated: %s", string(kv.Value))
}

// Delete
func handleDeploymentDelete(kv *mvccpb.KeyValue) {
	log.Printf("Deployment deleted: %s", string(kv.Value))

	var deployment shared.Deployment
	if err := json.Unmarshal(kv.Value, &deployment); err != nil {
		log.Printf("Failed to unmarshal deployment: %v", err)
		return
	}

	deletePodsByDeploymentID(deployment.ID)
}

func deletePodsByDeploymentID(deploymentID string) {
	pods, err := etcd.GetPodsByDeploymentID(deploymentID)
	if err != nil {
		log.Printf("Failed to get pods by deployment ID: %v", err)
		return
	}

	for _, pod := range pods {
		err := etcd.DeletePod(pod.ID)
		if err != nil {
			log.Printf("Failed to delete pod %s: %v", pod.ID, err)
			return
		}
	}
}