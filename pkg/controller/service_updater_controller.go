package controller

import (
	"maden/pkg/shared"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

type DefaultServiceUpdaterController struct {
}

func NewDefaultServiceUpdaterController() ServiceUpdaterController {
	return &DefaultServiceUpdaterController{}
}

// To be implemented once DNS server is set up
func (c *DefaultServiceUpdaterController) HandleServiceCreate(kv *mvccpb.KeyValue) {
	shared.Log.Infof("New service created: %s", string(kv.Value))
}

func (c *DefaultServiceUpdaterController) HandleServiceUpdate(prevKv *mvccpb.KeyValue, newKv *mvccpb.KeyValue) {
	shared.Log.Infof("Service updated: %s, %v", string(prevKv.Value), string(newKv.Value))
}

func (c *DefaultServiceUpdaterController) HandleServiceDelete(prevKv *mvccpb.KeyValue) {
	shared.Log.Infof("Service deleted: %s", string(prevKv.Value))
}
