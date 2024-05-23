package orchestrator

import (
	"maden/pkg/etcd"
	"maden/pkg/madelet"
	"maden/pkg/scheduler"
	"maden/pkg/shared"

)

type DefaultPodOrchestrator struct {
	Repo etcd.PodRepository
	Scheduler scheduler.Scheduler
}

func NewDefaultPodOrchestrator(
	repo etcd.PodRepository,
	scheduler scheduler.Scheduler,
) PodOrchestrator {
	return &DefaultPodOrchestrator{Repo: repo, Scheduler: scheduler}
}

func (po *DefaultPodOrchestrator) OrchestratePodCreation(pod *shared.Pod) error {
	err := po.Scheduler.SchedulePod(pod)
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