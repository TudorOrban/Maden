package controller

import (
	"maden/pkg/orchestrator"
	"maden/pkg/shared"
)

type DefaultPersistentVolumeClaimController struct {
	Orchestrator orchestrator.PersistentVolumeClaimOrchestrator
}

func NewDefaultPersistentVolumeClaimController(orchestrator orchestrator.PersistentVolumeClaimOrchestrator) PersistentVolumeClaimController {
	return &DefaultPersistentVolumeClaimController{Orchestrator: orchestrator}
}

func (po *DefaultPersistentVolumeClaimController) HandleIncomingPersistentVolumeClaim(volumeClaimSpec shared.PersistentVolumeClaimSpec) error {
	return po.Orchestrator.OrchestratePersistentVolumeClaimCreation(&volumeClaimSpec)
}
