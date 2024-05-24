package madelet

import (
	"fmt"
	"maden/pkg/shared"

	"context"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func NewDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	return cli
}

type DockerRuntime struct {
	Client *client.Client
}

func NewContainerRuntimeInterface(client *client.Client) ContainerRuntimeInterface {
	return &DockerRuntime{Client: client}	
}


func (d *DockerRuntime) CreateContainer(image string) (string, error) {
	log.Printf("Creating container with image %s", image)

	ctx := context.Background()
	resp, err := d.Client.ContainerCreate(ctx, &container.Config{
		Image: image,
	}, nil, nil, nil, "")
	if err != nil {
		log.Printf("Failed to create container: %v", err)
		return "", err
	}

	return resp.ID, nil
}

func (d *DockerRuntime) StartContainer(containerID string) error {
	log.Printf("Starting container %s", containerID)

	ctx := context.Background()
	if err := d.Client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		log.Printf("Failed to start container %s: %v", containerID, err)
		return err
	}

	return nil
}

func (d *DockerRuntime) StopContainer(containerID string) error {
	log.Printf("Stopping container %s", containerID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := d.Client.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		log.Printf("Failed to stop container %s: %v", containerID, err)
		return err
	}

	log.Printf("Container %s stopped", containerID)
	return nil
}

func (d *DockerRuntime) DeleteContainer(containerID string) error {
	log.Printf("Deleting container %s", containerID)

	ctx := context.Background()

	options := container.RemoveOptions{Force: true}

	if err := d.Client.ContainerRemove(ctx, containerID, options); err != nil {
		log.Printf("Failed to delete container %s: %v", containerID, err)
		return err
	}

	log.Printf("Container %s deleted", containerID)
	return nil
}

func (d *DockerRuntime) GetContainerLogs(containerID string, follow bool) (io.ReadCloser, error) {
	log.Printf("Getting logs for container %s", containerID)

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Hour) // Set a reasonable timeout
	// defer cancel()

	options := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: follow}
	return d.Client.ContainerLogs(ctx, containerID, options)
}

func (d *DockerRuntime) GetContainerStatus(containerID string) (shared.ContainerStatus, error) {
	if containerID == "" {
		log.Println("Empty container ID provided")
		return shared.Dead, fmt.Errorf("empty container ID")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	resp, err := d.Client.ContainerInspect(ctx, containerID)
	if err != nil {
		log.Printf("Failed to get status for container %s: %v", containerID, err)
		return shared.Dead, err
	}

	containerStatus, err := shared.GetStatusFromString(resp.State.Status)
	if err != nil {
		log.Printf("Failed to parse status for container %s: %v", containerID, err)
		return shared.Dead, err
	}
	return *containerStatus, nil
}
