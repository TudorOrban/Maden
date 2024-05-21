package etcd

import (
	"context"
	"encoding/json"
	"maden/pkg/shared"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var servicesKey = "services/"


func ListServices() ([]shared.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, servicesKey, clientv3.WithPrefix())
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

func GetServiceByName(name string) (*shared.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	key := servicesKey + name
	resp, err := cli.Get(ctx, key)
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

func CreateService(service *shared.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    serviceData, err := json.Marshal(service)
    if err != nil {
        return err
    }

	key := servicesKey + service.Name

	txnResp, err := cli.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, string(serviceData))).
		Else(clientv3.OpGet(key)).
		Commit()

	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return &shared.ErrDuplicateResource{ID: service.ID, ResourceType: shared.ServiceResource}
	}

	return nil
}

func UpdateService(service *shared.Service) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    serviceData, err := json.Marshal(service)
    if err != nil {
        return err
    }

    key := servicesKey + service.Name

    _, err = cli.Put(ctx, key, string(serviceData))
    if err != nil {
        return err
    }

    return nil
}

func DeleteService(serviceID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := servicesKey + serviceID

	resp, err := cli.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return &shared.ErrNotFound{ID: serviceID, ResourceType: shared.ServiceResource}
	}
	return nil
}