package madelet

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"
	"time"

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
		pod.Status = shared.PodContainerCreating
		if err := p.PodRepo.UpdatePod(pod); err != nil {
			log.Printf("Failed to update pod status: %v", err)
			return
		}
		time.Sleep(5 * time.Second) // Sleep for 5 seconds
		
		containerID, err := p.Runtime.CreateContainer(container.Image)
		if err != nil {
			log.Printf("Failed to create container: %v", err)
			return
		}

		if err := p.Runtime.StartContainer(containerID); err != nil {
			log.Printf("Failed to start container: %v", err)
			return
		}
		pod.Status = shared.PodRunning
		if err := p.PodRepo.UpdatePod(pod); err != nil {
			log.Printf("Failed to update pod status: %v", err)
			return
		}

		if err := p.Runtime.GetContainerLogs(containerID, true); err != nil {
			log.Printf("Failed to get container logs: %v", err)
			return
		}
	}
}