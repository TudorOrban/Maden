package madelet

import (
	"context"
	"errors"
	"maden/pkg/mocks"
	"maden/pkg/shared"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateContainer(t *testing.T) {

	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	
	mockClient := mocks.NewMockDockerClient(ctrl)
	runtime := NewContainerRuntimeInterface(mockClient)

	ctx := context.Background()
	containerID := "abc123"
	image := "nginx:latest"

	mockClient.EXPECT().ContainerCreate(ctx, gomock.Any(), nil, nil, nil, "").Return(container.CreateResponse{ID: containerID}, nil)

	// Act
	id, err := runtime.CreateContainer(image)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, containerID, id)

	mockClient.EXPECT().
		ContainerCreate(ctx, gomock.Any(), nil, nil, nil, "").
		Return(container.CreateResponse{}, errors.New("error creating container"))

	// Act
	_, err = runtime.CreateContainer(image)

	// Assert
	assert.Error(t, err)
}

func TestStartContainer(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockDockerClient(ctrl)
	runtime := NewContainerRuntimeInterface(mockClient)

	ctx := context.Background()
	containerID := "abc123"

	// Success expectation
	mockClient.EXPECT().ContainerStart(ctx, containerID, container.StartOptions{}).Return(nil)

	// Act
	err := runtime.StartContainer(containerID)

	// Assert
	assert.NoError(t, err)

	// Failure expectation
	mockClient.EXPECT().ContainerStart(ctx, containerID, container.StartOptions{}).
		Return(errors.New("error starting container"))

	// Act
	err = runtime.StartContainer(containerID)

	// Assert
	assert.Error(t, err)
}

func TestStopContainer(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockDockerClient(ctrl)
    runtime := NewContainerRuntimeInterface(mockClient)

    containerID := "abc123"

    // Success expectation
    mockClient.EXPECT().ContainerStop(gomock.Any(), containerID, container.StopOptions{}).Return(nil)

    // Act
    err := runtime.StopContainer(containerID)

    // Assert
    assert.NoError(t, err)

    // Failure expectation
    mockClient.EXPECT().ContainerStop(gomock.Any(), containerID, container.StopOptions{}).
        Return(errors.New("error stopping container"))

    // Act
    err = runtime.StopContainer(containerID)

    // Assert
    assert.Error(t, err)
}

func TestDeleteContainer(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockDockerClient(ctrl)
	runtime := NewContainerRuntimeInterface(mockClient)

	ctx := context.Background()
	containerID := "abc123"

	// Success expectation
	mockClient.EXPECT().ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}).Return(nil)

	// Act
	err := runtime.DeleteContainer(containerID)

	// Assert
	assert.NoError(t, err)

	// Failure expectation
	mockClient.EXPECT().ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true}).
		Return(errors.New("error deleting container"))

	// Act
	err = runtime.DeleteContainer(containerID)

	// Assert
	assert.Error(t, err)
}

func TestGetContainerStatus(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockDockerClient(ctrl)
    runtime := NewContainerRuntimeInterface(mockClient)

    containerID := "abc123"
    status := "running"
    expectedStatus, _ := shared.GetStatusFromString(status)

    // Mocking the successful inspect call
    mockClient.EXPECT().ContainerInspect(gomock.Any(), containerID).Return(types.ContainerJSON{
        ContainerJSONBase: &types.ContainerJSONBase{
            State: &types.ContainerState{
                Status: status,
            },
        },
    }, nil)

    // Act
    actualStatus, err := runtime.GetContainerStatus(containerID)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, *expectedStatus, actualStatus)

    // Mocking a failure in the inspect call
    mockClient.EXPECT().ContainerInspect(gomock.Any(), containerID).
        Return(types.ContainerJSON{}, errors.New("container not found"))

    // Act
    _, err = runtime.GetContainerStatus(containerID)

    // Assert
    assert.Error(t, err)
}
