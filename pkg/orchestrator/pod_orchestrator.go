package orchestrator

import (
	"maden/pkg/etcd"
	"maden/pkg/madelet"
	"maden/pkg/scheduler"
	"maden/pkg/shared" 

	"context"
	"fmt"
	"io"
)

type DefaultPodOrchestrator struct {
	Repo etcd.PodRepository
	Scheduler scheduler.Scheduler
	PodManager madelet.PodManager
}

func NewDefaultPodOrchestrator(
	repo etcd.PodRepository,
	scheduler scheduler.Scheduler,
	podManager madelet.PodManager,
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

func (po *DefaultPodOrchestrator) GetPodLogs(ctx context.Context, podID string, containerID string, follow bool) (io.ReadCloser, error) {
	pod, err := po.Repo.GetPodByID(podID)
	if err != nil {
		return nil, err
	}

	actualContainerID, err := determineContainerID(pod, containerID)
	if err != nil {
		return nil, err
	}

	return po.PodManager.GetContainerLogs(ctx, actualContainerID, follow)
}

func (po *DefaultPodOrchestrator) OrchestrateContainerCommandExecution(ctx context.Context, podID string, containerID string, cmd string) (string, error) {
	pod, err := po.Repo.GetPodByID(podID)
	if err != nil {
		return "", err
	}

	actualContainerID, err := determineContainerID(pod, containerID)
	if err != nil {
		return "", err
	}

	return po.PodManager.ExecuteCommandInContainer(ctx, actualContainerID, cmd)
}

func determineContainerID(pod *shared.Pod, containerID string) (string, error) {
	if len(pod.Containers) > 1 {
		if containerID == "" {
			return "", fmt.Errorf("containerID is required for pods with multiple containers")
		}
		return containerID, nil
	}
	return pod.Containers[0].ID, nil
}