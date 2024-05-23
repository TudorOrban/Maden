package controller

import (
	"maden/pkg/shared"

	"encoding/json"
	"log"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

type DefaultPodUpdaterController struct {

}

func NewDefaultPodUpdaterController() PodUpdaterController {
	return &DefaultPodUpdaterController{}
}


func (c *DefaultPodUpdaterController) HandlePodUpdate(oldKv *mvccpb.KeyValue, newKv *mvccpb.KeyValue) {
	log.Printf("Pod updated: %s, %v", string(oldKv.Value), string(newKv.Value))
	
	var newPod shared.Pod
	if err := json.Unmarshal(newKv.Value, &newPod); err != nil {
		log.Printf("Failed to unmarshal new pod: %v", err)
		return
	}
}