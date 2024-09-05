package madelet

import (
	"maden/pkg/shared"

	"context"
	"fmt"
	"io"
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
		shared.Log.Errorf("Failed to create Docker client: %v", err)
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
	ctx := context.Background()
	resp, err := d.Client.ContainerCreate(ctx, &container.Config{
		Image: image,
	}, nil, nil, nil, "")
	if err != nil {
		shared.Log.Errorf("Failed to create container: %v", err)
		return "", err
	}

	return resp.ID, nil
}

func (d *DockerRuntime) StartContainer(containerID string) error {
	ctx := context.Background()
	if err := d.Client.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		shared.Log.Errorf("Failed to start container %s: %v", containerID, err)
		return err
	}

	shared.Log.Infof("Container %s started", containerID)
	return nil
}

func (d *DockerRuntime) StopContainer(containerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := d.Client.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		shared.Log.Errorf("Failed to stop container %s: %v", containerID, err)
		return err
	}

	shared.Log.Infof("Container %s stopped", containerID)
	return nil
}

func (d *DockerRuntime) DeleteContainer(containerID string) error {
	ctx := context.Background()

	options := container.RemoveOptions{Force: true}

	if err := d.Client.ContainerRemove(ctx, containerID, options); err != nil {
		shared.Log.Errorf("Failed to delete container %s: %v", containerID, err)
		return err
	}

	shared.Log.Infof("Container %s deleted", containerID)
	return nil
}

func (d *DockerRuntime) GetContainerLogs(ctx context.Context, containerID string, follow bool) (io.ReadCloser, error) {
	shared.Log.Infof("Getting logs for container %s", containerID)

	options := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: follow}
	return d.Client.ContainerLogs(ctx, containerID, options)
}

func (d *DockerRuntime) GetContainerStatus(containerID string) (shared.ContainerStatus, error) {
	if containerID == "" {
		shared.Log.Errorf("Empty container ID provided")
		return shared.Dead, fmt.Errorf("empty container ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := d.Client.ContainerInspect(ctx, containerID)
	if err != nil {
		shared.Log.Errorf("Failed to get status for container %s: %v", containerID, err)
		return shared.Dead, err
	}

	containerStatus, err := shared.GetStatusFromString(resp.State.Status)
	if err != nil {
		shared.Log.Errorf("Failed to parse status for container %s: %v", containerID, err)
		return shared.Dead, err
	}
	return *containerStatus, nil
}

func (d *DockerRuntime) ExecCommandCreate(ctx context.Context, containerID string, execConfig types.ExecConfig) (string, error) {
	execID, err := d.Client.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		shared.Log.Errorf("Failed to create exec command: %v", err)
		return "", err
	}

	return execID.ID, nil
}

func (d *DockerRuntime) ExecCommandAttach(ctx context.Context, execID string, attachConfig types.ExecStartCheck, tty bool) (*types.HijackedResponse, error) {
	execAttach, err := d.Client.ContainerExecAttach(ctx, execID, attachConfig)
	if err != nil {
		shared.Log.Errorf("Failed to attach to exec command: %v", err)
		return nil, err
	}

	return &execAttach, nil
}
