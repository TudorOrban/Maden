package madelet

import (
	"io"
	"maden/pkg/shared"
)

type ContainerRuntimeInterface interface {
	CreateContainer(image string) (string, error)
	StartContainer(containerID string) error
	StopContainer(containerID string) error
	DeleteContainer(containerID string) error
	GetContainerLogs(containerID string, follow bool) (io.ReadCloser, error)
	GetContainerStatus(containerID string) (shared.ContainerStatus, error)
}