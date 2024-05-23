package etcd

import (
	"context"
	"encoding/json"
	"maden/pkg/shared"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var servicesKey = "services/"

type EtcdServiceRepository struct {
	client EtcdClient
	transactioner Transactioner
}

func NewEtcdServiceRepository(
	client EtcdClient,
	transactioner Transactioner,
) ServiceRepository {
	return &EtcdServiceRepository{client: client, transactioner: transactioner}
}


func (repo *EtcdServiceRepository) ListServices() ([]shared.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := repo.client.Get(ctx, servicesKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	services := make([]shared.Service, 0)
	for _, kv := range resp.Kvs {
		var service shared.Service
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

func (repo *EtcdServiceRepository) GetServiceByName(name string) (*shared.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	key := servicesKey + name
	resp, err := repo.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, &shared.ErrNotFound{Name: name, ResourceType: shared.ServiceResource}
	}

	var service shared.Service
	if err := json.Unmarshal(resp.Kvs[0].Value, &service); err != nil {
		return nil, err
	}
	return &service, nil
} 

func (repo *EtcdServiceRepository) CreateService(service *shared.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    serviceData, err := json.Marshal(service)
    if err != nil {
        return err
    }

	key := servicesKey + service.Name

	return repo.transactioner.PerformTransaction(ctx, key, string(serviceData), shared.ServiceResource)
}

func (repo *EtcdServiceRepository) UpdateService(service *shared.Service) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    serviceData, err := json.Marshal(service)
    if err != nil {
        return err
    }

    key := servicesKey + service.Name

    resp, err := repo.client.Put(ctx, key, string(serviceData), clientv3.WithPrevKV())
    if err != nil {
        return err
    }

	if resp.PrevKv == nil {
		return &shared.ErrNotFound{ID: service.ID, ResourceType: shared.ServiceResource}
	}
	return nil
}

func (repo *EtcdServiceRepository) DeleteService(serviceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := servicesKey + serviceName

	resp, err := repo.client.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return &shared.ErrNotFound{ID: serviceName, ResourceType: shared.ServiceResource}
	}
	return nil
}