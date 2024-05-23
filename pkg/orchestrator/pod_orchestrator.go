package orchestrator

import (
	"maden/pkg/etcd"
	"maden/pkg/madelet"
	"maden/pkg/scheduler"
	"maden/pkg/shared"

)

type DefaultPodOrchestrator struct {
	Repo etcd.PodRepository
}

func NewDefaultPodOrchestrator(repo etcd.PodRepository) PodOrchestrator {
	return &DefaultPodOrchestrator{Repo: repo}
}

func (po *DefaultPodOrchestrator) OrchestratePodCreation(pod *shared.Pod) error {
	err := scheduler.SchedulePod(pod)
	if err != nil {
		return err
	}

	if err := po.Repo.CreatePod(pod); err != nil {
		return err
	}

	dockerRuntime, err := madelet.NewDockerRuntime()
	if err != nil {
		return err
	}

	podLifecycleManager := madelet.PodLifecycleManager{
		Runtime: dockerRuntime,
	}
	go podLifecycleManager.RunPod(pod)

	return nil
}