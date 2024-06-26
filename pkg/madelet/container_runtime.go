package madelet

import (
	"maden/pkg/shared"

	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func NewDockerClient(client *client.Client) DockerClient {
	return client
}

func NewClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	return cli
}

type DockerRuntime struct {
	Client DockerClient
}

func NewContainerRuntimeInterface(client DockerClient) ContainerRuntimeInterface {
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

func (d *DockerRuntime) GetContainerLogs(ctx context.Context, containerID string, follow bool) (io.ReadCloser, error) {
	log.Printf("Getting logs for container %s", containerID)

	options := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: follow}
	return d.Client.ContainerLogs(ctx, containerID, options)
}

func (d *DockerRuntime) GetContainerStatus(containerID string) (shared.ContainerStatus, error) {
	if containerID == "" {
		log.Println("Empty container ID provided")
		return shared.Dead, fmt.Errorf("empty container ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

func (d *DockerRuntime) ExecCommandCreate(ctx context.Context, containerID string, execConfig types.ExecConfig) (string, error) {
	execID, err := d.Client.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		log.Printf("Failed to create exec command: %v", err)
		return "", err
	}

	return execID.ID, nil
}

func (d *DockerRuntime) ExecCommandAttach(ctx context.Context, execID string, attachConfig types.ExecStartCheck, tty bool) (*types.HijackedResponse, error) {
	execAttach, err := d.Client.ContainerExecAttach(ctx, execID, attachConfig)
	if err != nil {
		log.Printf("Failed to attach to exec command: %v", err)
		return nil, err
	}

	return &execAttach, nil
}