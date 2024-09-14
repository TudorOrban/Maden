package orchestrator

import (
	"context"
	"io"
	"maden/pkg/shared"
)

type PodOrchestrator interface {
	OrchestratePodCreation(pod *shared.Pod) error
	OrchestratePodDeletion(pod *shared.Pod) error
	GetPodLogs(ctx context.Context, podID string, containerID string, follow bool) (io.ReadCloser, error)
	OrchestrateContainerCommandExecution(ctx context.Context, podID string, containerID string, cmd string) (string, error)
}

type ServiceOrchestrator interface {
	OrchestrateServiceCreation(serviceSpec shared.ServiceSpec) error
	OrchestrateServiceUpdate(existingService shared.Service, serviceSpec shared.ServiceSpec) error
	OrchestrateServiceDeletion(serviceName string) error
}

type PersistentVolumeOrchestrator interface {
	OrchestratePersistentVolumeCreation(volumeSpec *shared.PersistentVolumeSpec) error
	OrchestratePersistentVolumeUpdate(existingVolume *shared.PersistentVolume, volumeSpec *shared.PersistentVolumeSpec) error
	OrchestratePersistentVolumeDeletion(volumeName string) error
}

type PersistentVolumeClaimOrchestrator interface {
	OrchestratePersistentVolumeClaimCreation(volumeClaimSpec *shared.PersistentVolumeClaimSpec) error
	OrchestratePersistentVolumeClaimUpdate(existingVolumeClaim *shared.PersistentVolumeClaim, volumeClaimSpec *shared.PersistentVolumeClaimSpec) error
	OrchestratePersistentVolumeClaimDeletion(volumeClaimName string) error
}
