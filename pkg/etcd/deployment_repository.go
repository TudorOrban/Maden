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

	resp, err := cli.Get(ctx, deploymentsKey, clientv3.WithPrefix())
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

func CreateDeployment(deployment *shared.Deployment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    deploymentData, err := json.Marshal(deployment)
    if err != nil {
        return err
    }

	key := deploymentsKey + deployment.Name

	txnResp, err := cli.Txn(ctx).
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

func DeleteDeployment(deploymentID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := deploymentsKey + deploymentID

	resp, err := cli.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return shared.ErrNotFound
	}
	return nil
}