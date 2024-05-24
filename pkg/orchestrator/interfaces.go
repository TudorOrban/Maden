package orchestrator

import (
	"io"
	"maden/pkg/shared"
)

type PodOrchestrator interface {
	OrchestratePodCreation(pod *shared.Pod) error
	OrchestratePodDeletion(pod *shared.Pod) error
	GetPodLogs(podID string, containerID string, follow bool) (io.ReadCloser, error)
}