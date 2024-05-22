package madelet

import "log"

type ContainerRuntimeInterface interface {
	CreateContainer(image string) (string, error)
	StartContainer(containerID string) error
	StopContainer(containerID string) error
	DeleteContainer(containerID string) error
}

type DockerRuntime struct{}

func (d *DockerRuntime) CreateContainer(image string) (string, error) {
	log.Printf("Creating container with image %s", image)
	return "docker-container-id", nil
}

func (d *DockerRuntime) StartContainer(containerID string) error {
	log.Printf("Starting container %s", containerID)
	return nil
}

func (d *DockerRuntime) StopContainer(containerID string) error {
	log.Printf("Stopping container %s", containerID)
	return nil
}

func (d *DockerRuntime) DeleteContainer(containerID string) error {
	log.Printf("Deleting container %s", containerID)
	return nil
}