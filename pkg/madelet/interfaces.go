package madelet

import (
	"context"
	"io"
	"maden/pkg/shared"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type DockerClient interface {
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error)
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error)
	ContainerExecCreate(ctx context.Context, containerID string, config types.ExecConfig) (types.IDResponse, error)
	ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
}

type ContainerRuntimeInterface interface {
	CreateContainer(image string) (string, error)
	StartContainer(containerID string) error
	StopContainer(containerID string) error
	DeleteContainer(containerID string) error
	GetContainerLogs(ctx context.Context, containerID string, follow bool) (io.ReadCloser, error)
	GetContainerStatus(containerID string) (shared.ContainerStatus, error)
}

type PodManager interface {
	RunPod(pod *shared.Pod)
	StopPod(pod *shared.Pod) error
	GetContainerLogs(ctx context.Context, containerID string, follow bool) (io.ReadCloser, error)
	ExecuteCommandInContainer(containerID string, command string) (string, error)
}