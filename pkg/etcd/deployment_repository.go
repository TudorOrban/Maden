package etcd

import (
	"context"
	"encoding/json"
	"maden/pkg/shared"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var deploymentsKey = "deployments/"


func ListDeployments() ([]shared.Deployment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := Cli.Get(ctx, deploymentsKey, clientv3.WithPrefix())
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

func GetDeploymentByName(name string) (*shared.Deployment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

	key := deploymentsKey + name
	resp, err := Cli.Get(ctx, key)
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

func CreateDeployment(deployment *shared.Deployment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    deploymentData, err := json.Marshal(deployment)
    if err != nil {
        return err
    }

	key := deploymentsKey + deployment.Name

	txnResp, err := Cli.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, string(deploymentData))).
		Else(clientv3.OpGet(key)).
		Commit()

	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return &shared.ErrDuplicateResource{ID: deployment.ID, ResourceType: shared.DeploymentResource}
	}

	return nil
}

func UpdateDeployment(deployment *shared.Deployment) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    deploymentData, err := json.Marshal(deployment)
    if err != nil {
        return err
    }

    key := deploymentsKey + deployment.Name

    _, err = Cli.Put(ctx, key, string(deploymentData))
    if err != nil {
        return err
    }

    return nil
}

func DeleteDeployment(deploymentName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := deploymentsKey + deploymentName

	resp, err := Cli.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return &shared.ErrNotFound{ID: deploymentName, ResourceType: shared.DeploymentResource}
	}
	return nil
}