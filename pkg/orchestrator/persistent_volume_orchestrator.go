package orchestrator

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"
)

type DefaultPersistentVolumeOrchestrator struct {
	Repo etcd.PersistentVolumeRepository
}

func NewDefaultPersistentVolumeOrchestrator(repo etcd.PersistentVolumeRepository) PersistentVolumeOrchestrator {
	return &DefaultPersistentVolumeOrchestrator{Repo: repo}
}

func (po *DefaultPersistentVolumeOrchestrator) OrchestratePersistentVolumeCreation(volumeSpec *shared.PersistentVolumeSpec) error {
	var volumeID = volumeSpec.Name + shared.GenerateRandomString(10)

	var volume = &shared.PersistentVolume{
		ID:                            volumeID,
		Name:                          volumeSpec.Name,
		Capacity:                      volumeSpec.Capacity,
		AccessModes:                   volumeSpec.AccessModes,
		PersistentVolumeReclaimPolicy: volumeSpec.PersistentVolumeReclaimPolicy,
		StorageClassName:              volumeSpec.StorageClassName,
		MountOptions:                  volumeSpec.MountOptions,
	}

	return po.Repo.CreatePersistentVolume(volume)
}

func (po *DefaultPersistentVolumeOrchestrator) OrchestratePersistentVolumeUpdate(existingVolume *shared.PersistentVolume, volumeSpec *shared.PersistentVolumeSpec) error {
	var volumeID = volumeSpec.Name + shared.GenerateRandomString(10)

	var volume = &shared.PersistentVolume{
		ID:                            volumeID,
		Name:                          volumeSpec.Name,
		Capacity:                      volumeSpec.Capacity,
		AccessModes:                   volumeSpec.AccessModes,
		PersistentVolumeReclaimPolicy: volumeSpec.PersistentVolumeReclaimPolicy,
		StorageClassName:              volumeSpec.StorageClassName,
		MountOptions:                  volumeSpec.MountOptions,
	}

	return po.Repo.UpdatePersistentVolume(volume)
}

func (po *DefaultPersistentVolumeOrchestrator) OrchestratePersistentVolumeDeletion(pvID string) error {
	return po.Repo.DeletePersistentVolume(pvID)
}
