package controller

import (
	"log"

	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdChangeListener struct {
	client *clientv3.Client
	DeploymentController DeploymentUpdaterController
	ServiceController ServiceUpdaterController
}

func NewEtcdChangeListener(
	client *clientv3.Client,
	deploymentController DeploymentUpdaterController,
	serviceController ServiceUpdaterController,
) *EtcdChangeListener {
	return &EtcdChangeListener{client: client, DeploymentController: deploymentController, ServiceController: serviceController}
}

func (l *EtcdChangeListener) WatchDeployments() {	
	ctx := context.Background()
	rch := l.client.Watch(ctx, "deployments/", clientv3.WithPrefix(), clientv3.WithPrevKV())
	log.Println("Watching deployments...")

	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					l.DeploymentController.HandleDeploymentCreate(ev.Kv)
				} else {
					l.DeploymentController.HandleDeploymentUpdate(ev.PrevKv, ev.Kv)
				}
			case clientv3.EventTypeDelete:
				l.DeploymentController.HandleDeploymentDelete(ev.PrevKv)
			}
		}
	}
}

func (l *EtcdChangeListener) WatchServices() {
	ctx := context.Background()
	rch := l.client.Watch(ctx, "services/", clientv3.WithPrefix(), clientv3.WithPrevKV())
	log.Println("Watching services...")

	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					l.ServiceController.HandleServiceCreate(ev.Kv)
				} else {
					l.ServiceController.HandleServiceUpdate(ev.PrevKv, ev.Kv)
				}
			case clientv3.EventTypeDelete:
				l.ServiceController.HandleServiceDelete(ev.PrevKv)
			}
		}
	}
}