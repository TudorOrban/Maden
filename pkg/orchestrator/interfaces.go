package orchestrator

import "maden/pkg/shared"

type PodOrchestrator interface {
	OrchestratePodCreation(pod *shared.Pod) error
}