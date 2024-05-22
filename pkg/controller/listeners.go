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
			log.Println("Key:", string(ev.Kv.Key))
			log.Println("Value:", string(ev.Kv.Value))
			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					handleDeploymentCreate(ev.Kv)
				} else {
					handleDeploymentUpdate(ev.Kv)
				}
			case clientv3.EventTypeDelete:
				log.Println("Deployment deleted in switch")
				handleDeploymentDelete(ev.PrevKv)
			}
		}
	}
}