package madelet

import (
	"log"
	"maden/pkg/shared"
)

type PodLifecycleManager struct {
	Runtime ContainerRuntimeInterface
}

func (p *PodLifecycleManager) RunPod(pod *shared.Pod) {
	for _, container := range pod.Containers {
		containerID, err := p.Runtime.CreateContainer(container.Image)
		if err != nil {
			log.Printf("Failed to create container: %v", err)
			return
		}
		if err := p.Runtime.StartContainer(containerID); err != nil {
			log.Printf("Failed to start container: %v", err)
			return
		}
		if err := p.Runtime.GetContainerLogs(containerID, true); err != nil {
			log.Printf("Failed to get container logs: %v", err)
			return
		}
	}
}