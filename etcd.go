package main

import (
	"log"
	"time"
	"go.etcd.io/etcd/client/v3"
)

var cli *clientv3.Client

func initEtcd() {
	var err error

	cli, err = clientv3.New(clientv3.Config{
		Endpoints: []string{"etcd:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
}