package orchestrator

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"
)

type DefaultPersistentVolumeClaimOrchestrator struct {
	Repo etcd.PersistentVolumeClaimRepository
}

func NewDefaultPersistentVolumeClaimOrchestrator(repo etcd.PersistentVolumeClaimRepository) PersistentVolumeClaimOrchestrator {
	return &DefaultPersistentVolumeClaimOrchestrator{Repo: repo}
}

func (po *DefaultPersistentVolumeClaimOrchestrator) OrchestratePersistentVolumeClaimCreation(volumeClaimSpec *shared.PersistentVolumeClaimSpec) error {
	var volumeClaimID = volumeClaimSpec.Name + shared.GenerateRandomString(10)

	var volumeClaim = &shared.PersistentVolumeClaim{
		ID: volumeClaimID,
		Name: volumeClaimSpec.Name,
		AccessModes: volumeClaimSpec.AccessModes,
		Resources: volumeClaimSpec.Resources,
		VolumeName: volumeClaimSpec.VolumeName,
	};

	return po.Repo.CreatePersistentVolumeClaim(volumeClaim)
}

func (po *DefaultPersistentVolumeClaimOrchestrator) OrchestratePersistentVolumeClaimUpdate(existingVolumeClaim *shared.PersistentVolumeClaim, volumeClaimSpec *shared.PersistentVolumeClaimSpec) error {
	var volumeClaimID = volumeClaimSpec.Name + shared.GenerateRandomString(10)

	var volumeClaim = &shared.PersistentVolumeClaim{
		ID: volumeClaimID,
		Name: volumeClaimSpec.Name,
		AccessModes: volumeClaimSpec.AccessModes,
		Resources: volumeClaimSpec.Resources,
		VolumeName: volumeClaimSpec.VolumeName,
	};

	return po.Repo.UpdatePersistentVolumeClaim(volumeClaim)
}

func (po *DefaultPersistentVolumeClaimOrchestrator) OrchestratePersistentVolumeClaimDeletion(pvID string) error {
	return po.Repo.DeletePersistentVolumeClaim(pvID)
}
