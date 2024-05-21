package etcd

import (
	"maden/pkg/shared"

	"context"
	"encoding/json"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var podsKey = "pods/";

func CreatePod(pod *shared.Pod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	podData, err := json.Marshal(pod)
	if err != nil {
		return err
	}

	key := podsKey + pod.ID

	// Start transaction to prevent duplicates
	txnResp, err := Cli.Txn(ctx).
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

func ListPods() ([]shared.Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := Cli.Get(ctx, podsKey, clientv3.WithPrefix())
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

func DeletePod(podID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := podsKey + podID

	resp, err := Cli.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return &shared.ErrNotFound{ID: podID, ResourceType: shared.PodResource}
	}
	return nil
}