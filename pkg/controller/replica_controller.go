package controller

import (
	"log"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

func handleDeploymentCreate(kv *mvccpb.KeyValue) {
	log.Printf("New deployment created: %s", string(kv.Value))
}

func handleDeploymentUpdate(kv *mvccpb.KeyValue) {
	log.Printf("Deployment updated: %s", string(kv.Value))
}

func handleDeploymentDelete(kv *mvccpb.KeyValue) {
	log.Printf("Deployment deleted: %s", string(kv.Value))
}