package madelet

import (
	"context"
	"io"
	"log"
	"os"

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

func (d *DockerRuntime) GetContainerLogs(containerID string, follow bool) error {
	log.Printf("Getting logs for container %s", containerID)

	ctx := context.Background()
	options := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: follow}
	out, err := d.Client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		log.Printf("Failed to get logs for container %s: %v", containerID, err)
		return err
	}
	defer out.Close()

	if follow {
		_, err = io.Copy(os.Stdout, out)
		if err != nil {
			log.Printf("Failed to stream logs for container %s: %v", containerID, err)
			return err
		}
	} else {
		logContents, err := io.ReadAll(out)
		if err != nil {
			log.Printf("Failed to read logs for container %s: %v", containerID, err)
			return err
		}

		log.Printf("Logs for container %s: %s", containerID, logContents)
	}
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