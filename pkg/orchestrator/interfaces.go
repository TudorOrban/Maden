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
}