package etcd

import (
	"log"
	"maden/pkg/shared"

	"context"
	"encoding/json"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var podsKey = "pods/";

type EtcdPodRepository struct {
	client EtcdClient
	transactioner Transactioner
}

func NewEtcdPodRepository(
	client EtcdClient,
	transactioner Transactioner,	
) PodRepository {
	return &EtcdPodRepository{client: client, transactioner: transactioner}
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
	log.Printf("Pods: %v", resp.Kvs)

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

func (repo *EtcdPodRepository) GetPodByID(podID string) (*shared.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := podsKey + podID

	resp, err := repo.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, &shared.ErrNotFound{ID: podID, ResourceType: shared.PodResource}
	}

	var pod shared.Pod
	if err := json.Unmarshal(resp.Kvs[0].Value, &pod); err != nil {
		return nil, err
	}
	return &pod, nil
}

func (repo *EtcdPodRepository) CreatePod(pod *shared.Pod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	podData, err := json.Marshal(pod)
	if err != nil {
		return err
	}

	key := podsKey + pod.ID

    return repo.transactioner.PerformTransaction(ctx, key, string(podData), shared.PodResource)
}

func (repo *EtcdPodRepository) UpdatePod(pod *shared.Pod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	podData, err := json.Marshal(pod)
	if err != nil {
		return err
	}

	key := podsKey + pod.ID

	resp, err := repo.client.Put(ctx, key, string(podData), clientv3.WithPrevKV())
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