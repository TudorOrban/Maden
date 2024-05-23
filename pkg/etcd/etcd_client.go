package etcd

import (
	"context"
	"log"
	"time"

	"go.etcd.io/etcd/client/v3"
)

type EtcdClient interface {
    Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
    Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
    Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error)
    Txn(ctx context.Context) clientv3.Txn
}

func NewEtcdClient(client *clientv3.Client) EtcdClient {
	return client
}

func NewClientv3() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"etcd:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	return client
}