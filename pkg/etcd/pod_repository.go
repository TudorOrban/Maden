package etcd

import (
	"maden/pkg/shared"

	"context"
	"encoding/json"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var podsKey = "pods/";

type EtcdPodRepository struct {
	client *clientv3.Client
}

func NewEtcdPodRepository(client *clientv3.Client) PodRepository {
	return &EtcdPodRepository{client: client}
}


func (repo *EtcdPodRepository) ListPods() ([]shared.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := repo.client.Get(ctx, podsKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	pods := make([]shared.Pod, 0)
	for _, kv := range resp.Kvs {
		var pod shared.Pod
		if err := json.Unmarshal(kv.Value, &pod); err != nil {
			return nil, err
		}
		pods = append(pods, pod)
	}
	return pods, nil
}

func (repo *EtcdPodRepository) GetPodsByDeploymentID(deploymentID string) ([]shared.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := repo.client.Get(ctx, podsKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	pods := make([]shared.Pod, 0)
	for _, kv := range resp.Kvs {
		var pod shared.Pod
		if err := json.Unmarshal(kv.Value, &pod); err != nil {
			return nil, err
		}
		if pod.DeploymentID == deploymentID {
			pods = append(pods, pod)
		}
	}
	return pods, nil
}

func (repo *EtcdPodRepository) CreatePod(pod *shared.Pod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	podData, err := json.Marshal(pod)
	if err != nil {
		return err
	}

	key := podsKey + pod.ID

	// Start transaction to prevent duplicates
	txnResp, err := repo.client.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, string(podData))).
		Else(clientv3.OpGet(key)).
		Commit()

	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return &shared.ErrDuplicateResource{ID: pod.ID, ResourceType: shared.PodResource}
	}

	return nil
}

func (repo *EtcdPodRepository) UpdatePod(pod *shared.Pod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	podData, err := json.Marshal(pod)
	if err != nil {
		return err
	}

	key := podsKey + pod.ID

	resp, err := repo.client.Put(ctx, key, string(podData))
	if err != nil {
		return err
	}

	if resp.PrevKv == nil {
		return &shared.ErrNotFound{ID: pod.ID, ResourceType: shared.PodResource}
	}
	return nil
}

func (repo *EtcdPodRepository) DeletePod(podID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := podsKey + podID

	resp, err := repo.client.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return &shared.ErrNotFound{ID: podID, ResourceType: shared.PodResource}
	}
	return nil
}