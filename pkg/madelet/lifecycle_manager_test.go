package madelet

import (
	"maden/pkg/mocks"
	"maden/pkg/shared"

	"errors"
	"testing"

	// "github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)


func TestPodLifecycleManagerRunPodSuccess(t *testing.T) {
 // Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntimeInterface(ctrl)
	mockPodRepo := mocks.NewMockPodRepository(ctrl)
	manager := NewPodLifecycleManager(mockRuntime, mockPodRepo)

	pod := &shared.Pod{
		Containers: []shared.Container{
			{Image: "example-image"},
		},
		Status: shared.PodPending,
	}

	// Expectations
	mockPodRepo.EXPECT().UpdatePod(gomock.Any()).AnyTimes().Return(nil) // Called twice, for creating and running status updates
	mockRuntime.EXPECT().CreateContainer(gomock.Any()).Return("containerID", nil)
	mockRuntime.EXPECT().StartContainer(gomock.Any()).Return(nil)

	// Act
	manager.RunPod(pod)

	// Assert
	assert.Equal(t, shared.PodRunning, pod.Status)
}

func TestPodLifecycleManagerRunPodFailCreateContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRuntime := mocks.NewMockContainerRuntimeInterface(ctrl)
	mockPodRepo := mocks.NewMockPodRepository(ctrl)
	manager := NewPodLifecycleManager(mockRuntime, mockPodRepo)

	pod := &shared.Pod{
		Containers: []shared.Container{
			{Image: "example-image"},
		},
		Status: shared.PodPending,
	}

	// Expectations
	mockPodRepo.EXPECT().UpdatePod(gomock.Any()).AnyTimes().Return(nil) // Update to failed status
	mockRuntime.EXPECT().CreateContainer(gomock.Any()).Return("", errors.New("creation error"))

	// Act
	manager.RunPod(pod)

	// Assert
	assert.Equal(t, shared.PodFailed, pod.Status)
}

// func TestPodLifecycleManagerExecuteCommandInContainer(t *testing.T) {
//     ctrl := gomock.NewController(t)
//     defer ctrl.Finish()

//     mockRuntime := mocks.NewMockContainerRuntimeInterface(ctrl)
//     manager := NewPodLifecycleManager(mockRuntime, nil) 

//     ctx := context.Background()
//     containerID := "test-container"
//     command := "echo Hello"

//     execID := "exec-id"
//     content := strings.NewReader("Hello")
//     bufReader := bufio.NewReader(content)
//     hijackedResponse := &types.HijackedResponse{Reader: bufReader}

//     mockRuntime.EXPECT().ExecCommandCreate(ctx, containerID, gomock.Any()).Return(execID, nil)
//     mockRuntime.EXPECT().ExecCommandAttach(ctx, execID, gomock.Any(), true).Return(hijackedResponse, nil)

//     output, err := manager.ExecuteCommandInContainer(ctx, containerID, command)

//     assert.NoError(t, err)
//     assert.Equal(t, "Hello", output)
// }
