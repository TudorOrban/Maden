package controller

import (
	"log"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

// To be implemented once starting to work with Docker
func handleServiceCreate(kv *mvccpb.KeyValue) {
	log.Printf("New service created: %s", string(kv.Value))
}

func handleServiceUpdate(prevKv *mvccpb.KeyValue, newKv *mvccpb.KeyValue) {
	log.Printf("Service updated: %s, %v", string(prevKv.Value), string(newKv.Value))
}

func handleServiceDelete(prevKv *mvccpb.KeyValue) {
	log.Printf("Service deleted: %s", string(prevKv.Value))
}