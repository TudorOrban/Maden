package madelet

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"

	"log"
)

type PodLifecycleManager struct {
	Runtime ContainerRuntimeInterface
	PodRepo etcd.PodRepository
}

func NewPodLifecycleManager(
	runtime ContainerRuntimeInterface,
	podRepo etcd.PodRepository,	
) PodLifecycleManager {
	return PodLifecycleManager{Runtime: runtime, PodRepo: podRepo}
}


func (p *PodLifecycleManager) RunPod(pod *shared.Pod) {
	for _, container := range pod.Containers {
		containerID := p.attemptContainerCreation(pod, container)
		if containerID == nil {
			return
		}

		p.attemptContainerStart(*containerID, pod)

		if err := p.Runtime.GetContainerLogs(*containerID, true); err != nil {
			log.Printf("Failed to get container logs: %v", err)
			return
		}
	}
}

func (p *PodLifecycleManager) attemptContainerCreation(pod *shared.Pod, container shared.Container) *string {
	pod.Status = shared.PodContainerCreating
	if err := p.PodRepo.UpdatePod(pod); err != nil {
		log.Printf("Failed to update pod status: %v", err)
		return nil
	}
	
	containerID, err := p.Runtime.CreateContainer(container.Image)
	if err != nil {
		log.Printf("Failed to create container: %v", err)
		pod.Status = shared.PodFailed
		_ = p.PodRepo.UpdatePod(pod)	
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