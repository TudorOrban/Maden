package etcd

import (
	"context"
	"maden/pkg/shared"
	"time"
)

var dnsKey = "dns/"

type EtcdDNSRepository struct {
	client EtcdClient
}

func NewEtcdDNSRepository(client EtcdClient) DNSRepository {
	return &EtcdDNSRepository{client: client}
}

func (repo *EtcdDNSRepository) RegisterService(serviceName string, serviceIP string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := dnsKey + serviceName
	_, err := repo.client.Put(ctx, key, serviceIP)
	return err
}

func (repo *EtcdDNSRepository) DeregisterService(serviceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := dnsKey + serviceName
	_, err := repo.client.Delete(ctx, key)
	return err
}

func (repo *EtcdDNSRepository) ResolveService(serviceName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := dnsKey + serviceName
	resp, err := repo.client.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", &shared.ErrNotFound{ID: serviceName, ResourceType: shared.DNSResource}
	}

	return string(resp.Kvs[0].Value), nil
}