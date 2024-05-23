package controller

import (
	"maden/pkg/madelet"
	"maden/pkg/shared"

	"encoding/json"
	"log"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

type DefaultPodUpdaterController struct {
	PodLifecycleManager madelet.PodLifecycleManager
}

func NewDefaultPodUpdaterController(
	podLifecycleManager madelet.PodLifecycleManager,
) PodUpdaterController {
	return &DefaultPodUpdaterController{PodLifecycleManager: podLifecycleManager}
}


func (c *DefaultPodUpdaterController) HandlePodUpdate(oldKv *mvccpb.KeyValue, newKv *mvccpb.KeyValue) {
	log.Printf("Pod updated: %s, %v", string(oldKv.Value), string(newKv.Value))
	
	var oldPod shared.Pod
	if err := json.Unmarshal(oldKv.Value, &oldPod); err != nil {
		log.Printf("Failed to unmarshal old pod: %v", err)
		return
	}

	var newPod shared.Pod
	if err := json.Unmarshal(newKv.Value, &newPod); err != nil {
		log.Printf("Failed to unmarshal new pod: %v", err)
		return
	}

	shouldRestart := shouldRestart(oldPod, newPod)
	if !shouldRestart {
		return
	}

	log.Printf("Restarting pod: %s", newPod.ID)
	c.PodLifecycleManager.RunPod(&newPod);
}

func shouldRestart(oldPod shared.Pod, newPod shared.Pod) bool {
	if oldPod.Status == newPod.Status {
		return false// Only care about status changes for now
	}

	if newPod.RestartPolicy == shared.RestartNever {
		return false// No need to restart
	}

	if newPod.Status != shared.PodFailed {
		return false// Only restart failed pods
	}

	return true
}