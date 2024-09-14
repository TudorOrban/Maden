package controller

import (
	"maden/pkg/orchestrator"
	"maden/pkg/shared"
)

type DefaultPersistentVolumeController struct {
	Orchestrator orchestrator.PersistentVolumeOrchestrator
}

func NewDefaultPersistentVolumeController(orchestrator orchestrator.PersistentVolumeOrchestrator) PersistentVolumeController {
	return &DefaultPersistentVolumeController{Orchestrator: orchestrator}
}

func (po *DefaultPersistentVolumeController) HandleIncomingPersistentVolume(volumeSpec shared.PersistentVolumeSpec) error {
	return po.Orchestrator.OrchestratePersistentVolumeCreation(&volumeSpec)
}
