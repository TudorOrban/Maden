package controller

import (
	"log"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

type DefaultServiceUpdaterController struct {
}

func NewDefaultServiceUpdaterController() ServiceUpdaterController {
	return &DefaultServiceUpdaterController{}
}

// To be implemented once starting to work with Docker
func (c *DefaultServiceUpdaterController) HandleServiceCreate(kv *mvccpb.KeyValue) {
	log.Printf("New service created: %s", string(kv.Value))
}

func (c *DefaultServiceUpdaterController) HandleServiceUpdate(prevKv *mvccpb.KeyValue, newKv *mvccpb.KeyValue) {
	log.Printf("Service updated: %s, %v", string(prevKv.Value), string(newKv.Value))
}

func (c *DefaultServiceUpdaterController) HandleServiceDelete(prevKv *mvccpb.KeyValue) {
	log.Printf("Service deleted: %s", string(prevKv.Value))
}