package etcd

import (
	"log"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3"
)

func NewEtcdClient() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"etcd:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	return client
}

func ProvideEtcdClient() *clientv3.Client {
	return NewEtcdClient()
}

var Cli *clientv3.Client

var Mu sync.Mutex

func InitEtcd() {
	var err error

	Cli, err = clientv3.New(clientv3.Config{
		Endpoints: []string{"etcd:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
}