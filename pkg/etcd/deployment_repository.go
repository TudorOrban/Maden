package etcd

import (
	"context"
	"encoding/json"
	"maden/pkg/shared"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var deploymentsKey = "deployments/"

type EtcdDeploymentRepository struct {
	client EtcdClient
	transactioner Transactioner
}

func NewEtcdDeploymentRepository(
	client EtcdClient,
	transactioner Transactioner,	
) DeploymentRepository {
	return &EtcdDeploymentRepository{client: client, transactioner: transactioner}
}


func (repo *EtcdDeploymentRepository) ListDeployments() ([]shared.Deployment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := repo.client.Get(ctx, deploymentsKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	deployments := make([]shared.Deployment, 0)
	for _, kv := range resp.Kvs {
		var deployment shared.Deployment
		if err := json.Unmarshal(kv.Value, &deployment); err != nil {
			return nil, err
		}
		deployments = append(deployments, deployment)
	}
	return deployments, nil
}

func (repo *EtcdDeploymentRepository) GetDeploymentByName(name string) (*shared.Deployment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	key := deploymentsKey + name
	resp, err := repo.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, &shared.ErrNotFound{Name: name, ResourceType: shared.DeploymentResource}
	}

	var deployment shared.Deployment
	if err := json.Unmarshal(resp.Kvs[0].Value, &deployment); err != nil {
		return nil, err
	}
	return &deployment, nil
} 

func (repo *EtcdDeploymentRepository) CreateDeployment(deployment *shared.Deployment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    deploymentData, err := json.Marshal(deployment)
    if err != nil {
        return err
    }

	key := deploymentsKey + deployment.Name

	return repo.transactioner.PerformTransaction(ctx, key, string(deploymentData), shared.DeploymentResource)
}

func (repo *EtcdDeploymentRepository) UpdateDeployment(deployment *shared.Deployment) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    deploymentData, err := json.Marshal(deployment)
    if err != nil {
        return err
    }

    key := deploymentsKey + deployment.Name

    resp, err := repo.client.Put(ctx, key, string(deploymentData))
    if err != nil {
        return err
    }
	
	if resp.PrevKv == nil {
		return &shared.ErrNotFound{ID: deployment.Name, ResourceType: shared.DeploymentResource}
	}
    return nil
}

func (repo *EtcdDeploymentRepository) DeleteDeployment(deploymentName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := deploymentsKey + deploymentName

	resp, err := repo.client.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return &shared.ErrNotFound{ID: deploymentName, ResourceType: shared.DeploymentResource}
	}
	return nil
}