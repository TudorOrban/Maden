package madelet

import (
	"context"
	"io"
	"maden/pkg/etcd"
	"maden/pkg/shared"

	"log"

	"github.com/docker/docker/api/types"
)

type PodLifecycleManager struct {
	Runtime ContainerRuntimeInterface
	PodRepo etcd.PodRepository
}

func NewPodLifecycleManager(
	runtime ContainerRuntimeInterface,
	podRepo etcd.PodRepository,	
) PodManager {
	return &PodLifecycleManager{Runtime: runtime, PodRepo: podRepo}
}


func (p *PodLifecycleManager) RunPod(pod *shared.Pod) {
	for containerIndex := range pod.Containers {
		containerID := p.attemptContainerCreation(pod, containerIndex)
		if containerID == nil {
			return
		}

		p.attemptContainerStart(*containerID, pod)
	}
}

func (p *PodLifecycleManager) attemptContainerCreation(pod *shared.Pod, containerIndex int) *string {
	pod.Status = shared.PodContainerCreating
	if err := p.PodRepo.UpdatePod(pod); err != nil {
		log.Printf("Failed to update pod status: %v", err)
		return nil
	}
	
	containerID, err := p.Runtime.CreateContainer(pod.Containers[containerIndex].Image)
	if err != nil {
		log.Printf("Failed to create container: %v", err)
		pod.Status = shared.PodFailed
		_ = p.PodRepo.UpdatePod(pod)	
		return nil
	}

	pod.Containers[containerIndex].ID = containerID
	if err := p.PodRepo.UpdatePod(pod); err != nil {
		log.Printf("Failed to update pod with ContainerID: %v", err)
		return nil
	}

	return &containerID
}

func (p *PodLifecycleManager) attemptContainerStart(containerID string, pod *shared.Pod) {
	if err := p.Runtime.StartContainer(containerID); err != nil {
		pod.Status = shared.PodFailed
		_ = p.PodRepo.UpdatePod(pod) 
		log.Printf("Failed to start container: %v", err)
		return
	}
	pod.Status = shared.PodRunning
	if err := p.PodRepo.UpdatePod(pod); err != nil {
		log.Printf("Failed to update pod status: %v", err)
		return
	}
}

func (p *PodLifecycleManager) StopPod(pod *shared.Pod) error {
	for _, container := range pod.Containers {
		containerStatus, err := p.Runtime.GetContainerStatus(container.ID)
		if err != nil {
			log.Printf("Failed to get container status: %v", err)
			continue
		}
		if containerStatus != shared.Running {
			continue
		}

		if err := p.Runtime.StopContainer(container.ID); err != nil {
			log.Printf("Failed to stop container: %v", err)
		}

		if err := p.Runtime.DeleteContainer(container.ID); err != nil {
			log.Printf("Failed to remove container: %v", err)
		}
	}
	return nil
}

func (p *PodLifecycleManager) GetContainerLogs(ctx context.Context, containerID string, follow bool) (io.ReadCloser, error) {
	return p.Runtime.GetContainerLogs(ctx, containerID, follow)
}

func (p *PodLifecycleManager) ExecuteCommandInContainer(ctx context.Context, containerID string, command string) (string, error) {
	execConfig := types.ExecConfig{
		Cmd: []string{"/bin/sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
		Tty: true,
	}

	execID, err := p.Runtime.ExecCommandCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", err
	}

	attachOptions := types.ExecStartCheck{Tty: execConfig.Tty}
	execAttach, err := p.Runtime.ExecCommandAttach(ctx, execID, attachOptions, execConfig.Tty)
	if err != nil {
		return "", err
	}
	defer execAttach.Close()

	output, err := io.ReadAll(execAttach.Reader)
	if err != nil {
		log.Printf("Failed to read exec output: %v", err)
		return "", err
	}

	return string(output), nil
}