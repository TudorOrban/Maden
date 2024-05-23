package controller

import (
	"log"
	"maden/pkg/etcd"

	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdChangeListener struct {
	client *clientv3.Client
	DeploymentController DeploymentUpdaterController
}

func NewEtcdChangeListener(
	client *clientv3.Client,
	deploymentController DeploymentUpdaterController,
) *EtcdChangeListener {
	return &EtcdChangeListener{client: client, DeploymentController: deploymentController}
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

func WatchServices() {
	ctx := context.Background()
	rch := etcd.Cli.Watch(ctx, "services/", clientv3.WithPrefix(), clientv3.WithPrevKV())
	log.Println("Watching services...")

	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					handleServiceCreate(ev.Kv)
				} else {
					handleServiceUpdate(ev.PrevKv, ev.Kv)
				}
			case clientv3.EventTypeDelete:
				handleServiceDelete(ev.PrevKv)
			}
		}
	}
}