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
	PodManager madelet.PodLifecycleManager
}

func NewDefaultPodOrchestrator(
	repo etcd.PodRepository,
	scheduler scheduler.Scheduler,
	podManager madelet.PodLifecycleManager,
) PodOrchestrator {
	return &DefaultPodOrchestrator{Repo: repo, Scheduler: scheduler, PodManager: podManager}
}

func (po *DefaultPodOrchestrator) OrchestratePodCreation(pod *shared.Pod) error {
	err := po.Scheduler.SchedulePod(pod)
	if err != nil {
		return err
	}

	if err := po.Repo.CreatePod(pod); err != nil {
		return err
	}

	go po.PodManager.RunPod(pod)

	return nil
}

func (po *DefaultPodOrchestrator) OrchestratePodDeletion(pod *shared.Pod) error {
	if err := po.PodManager.StopPod(pod); err != nil {
		return err
	}

	if err := po.Repo.DeletePod(pod.ID); err != nil {
		return err
	}

	return nil
}