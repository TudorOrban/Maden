package controller

import (
	"log"
	"maden/pkg/etcd"

	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func WatchDeployments() {	
	ctx := context.Background()
	rch := etcd.Cli.Watch(ctx, "deployments/", clientv3.WithPrefix(), clientv3.WithPrevKV())
	log.Println("Watching deployments...")

	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					handleDeploymentCreate(ev.Kv)
				} else {
					handleDeploymentUpdate(ev.PrevKv, ev.Kv)
				}
			case clientv3.EventTypeDelete:
				handleDeploymentDelete(ev.PrevKv)
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