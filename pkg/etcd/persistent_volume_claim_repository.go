package etcd

import (
	"maden/pkg/shared"

	"context"
	"encoding/json"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var pvcsKey = "pvs/"

type EtcdPersistentVolumeClaimRepository struct {
	client        EtcdClient
	transactioner Transactioner
}

func NewEtcdPersistentVolumeClaimRepository(
	client EtcdClient,
	transactioner Transactioner,
) PersistentVolumeClaimRepository {
	return &EtcdPersistentVolumeClaimRepository{client: client, transactioner: transactioner}
}

func (repo *EtcdPersistentVolumeClaimRepository) ListPersistentVolumeClaims() ([]shared.PersistentVolumeClaim, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := repo.client.Get(ctx, pvcsKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	persistentVolumeClaims := make([]shared.PersistentVolumeClaim, 0)
	for _, kv := range resp.Kvs {
		var persistentVolumeClaim shared.PersistentVolumeClaim
		if err := json.Unmarshal(kv.Value, &persistentVolumeClaim); err != nil {
			return nil, err
		}
		persistentVolumeClaims = append(persistentVolumeClaims, persistentVolumeClaim)
	}
	return persistentVolumeClaims, nil
}

func (repo *EtcdPersistentVolumeClaimRepository) GetPersistentVolumeClaimByID(persistentVolumeClaimID string) (*shared.PersistentVolumeClaim, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := pvcsKey + persistentVolumeClaimID

	resp, err := repo.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, &shared.ErrNotFound{ID: persistentVolumeClaimID, ResourceType: shared.PersistentVolumeClaimResource}
	}

	var persistentVolumeClaim shared.PersistentVolumeClaim
	if err := json.Unmarshal(resp.Kvs[0].Value, &persistentVolumeClaim); err != nil {
		return nil, err
	}
	return &persistentVolumeClaim, nil
}

func (repo *EtcdPersistentVolumeClaimRepository) CreatePersistentVolumeClaim(persistentVolumeClaim *shared.PersistentVolumeClaim) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	persistentVolumeClaimData, err := json.Marshal(persistentVolumeClaim)
	if err != nil {
		return err
	}

	key := pvcsKey + persistentVolumeClaim.ID

	return repo.transactioner.PerformTransaction(ctx, key, string(persistentVolumeClaimData), shared.PersistentVolumeClaimResource)
}

func (repo *EtcdPersistentVolumeClaimRepository) UpdatePersistentVolumeClaim(persistentVolumeClaim *shared.PersistentVolumeClaim) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	persistentVolumeClaimData, err := json.Marshal(persistentVolumeClaim)
	if err != nil {
		return err
	}

	key := pvcsKey + persistentVolumeClaim.ID

	resp, err := repo.client.Put(ctx, key, string(persistentVolumeClaimData), clientv3.WithPrevKV())
	if err != nil {
		return err
	}

	if resp.PrevKv == nil {
		return &shared.ErrNotFound{ID: persistentVolumeClaim.ID, ResourceType: shared.PersistentVolumeClaimResource}
	}
	return nil
}

func (repo *EtcdPersistentVolumeClaimRepository) DeletePersistentVolumeClaim(persistentVolumeClaimID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := pvcsKey + persistentVolumeClaimID

	resp, err := repo.client.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return &shared.ErrNotFound{ID: persistentVolumeClaimID, ResourceType: shared.PersistentVolumeClaimResource}
	}
	return nil
}
