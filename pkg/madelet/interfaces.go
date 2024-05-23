package madelet

import "maden/pkg/shared"

type ContainerRuntimeInterface interface {
	CreateContainer(image string) (string, error)
	StartContainer(containerID string) error
	StopContainer(containerID string) error
	DeleteContainer(containerID string) error
	GetContainerLogs(containerID string, follow bool) error
	GetContainerStatus(containerID string) (shared.ContainerStatus, error)
}