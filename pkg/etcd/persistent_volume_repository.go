package etcd

import (
	"maden/pkg/shared"

	"context"
	"encoding/json"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var pvsKey = "pvs/"

type EtcdPersistentVolumeRepository struct {
	client        EtcdClient
	transactioner Transactioner
}

func NewEtcdPersistentVolumeRepository(
	client EtcdClient,
	transactioner Transactioner,
) PersistentVolumeRepository {
	return &EtcdPersistentVolumeRepository{client: client, transactioner: transactioner}
}

func (repo *EtcdPersistentVolumeRepository) ListPersistentVolumes() ([]shared.PersistentVolume, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := repo.client.Get(ctx, pvsKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	persistentVolumes := make([]shared.PersistentVolume, 0)
	for _, kv := range resp.Kvs {
		var persistentVolume shared.PersistentVolume
		if err := json.Unmarshal(kv.Value, &persistentVolume); err != nil {
			return nil, err
		}
		persistentVolumes = append(persistentVolumes, persistentVolume)
	}
	return persistentVolumes, nil
}

func (repo *EtcdPersistentVolumeRepository) GetPersistentVolumeByID(persistentVolumeID string) (*shared.PersistentVolume, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := pvsKey + persistentVolumeID

	resp, err := repo.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, &shared.ErrNotFound{ID: persistentVolumeID, ResourceType: shared.PersistentVolumeResource}
	}

	var persistentVolume shared.PersistentVolume
	if err := json.Unmarshal(resp.Kvs[0].Value, &persistentVolume); err != nil {
		return nil, err
	}
	return &persistentVolume, nil
}

func (repo *EtcdPersistentVolumeRepository) CreatePersistentVolume(persistentVolume *shared.PersistentVolume) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	persistentVolumeData, err := json.Marshal(persistentVolume)
	if err != nil {
		return err
	}

	key := pvsKey + persistentVolume.ID

	return repo.transactioner.PerformTransaction(ctx, key, string(persistentVolumeData), shared.PersistentVolumeResource)
}

func (repo *EtcdPersistentVolumeRepository) UpdatePersistentVolume(persistentVolume *shared.PersistentVolume) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	persistentVolumeData, err := json.Marshal(persistentVolume)
	if err != nil {
		return err
	}

	key := pvsKey + persistentVolume.ID

	resp, err := repo.client.Put(ctx, key, string(persistentVolumeData), clientv3.WithPrevKV())
	if err != nil {
		return err
	}

	if resp.PrevKv == nil {
		return &shared.ErrNotFound{ID: persistentVolume.ID, ResourceType: shared.PersistentVolumeResource}
	}
	return nil
}

func (repo *EtcdPersistentVolumeRepository) DeletePersistentVolume(persistentVolumeID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := pvsKey + persistentVolumeID

	resp, err := repo.client.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return &shared.ErrNotFound{ID: persistentVolumeID, ResourceType: shared.PersistentVolumeResource}
	}
	return nil
}
