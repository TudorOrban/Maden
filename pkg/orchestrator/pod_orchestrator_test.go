package orchestrator

import (
	"errors"
	"maden/pkg/mocks"
	"maden/pkg/shared"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)


func TestOrchestratePodCreation(t *testing.T) {
	// Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockPodRepository(ctrl)
    mockScheduler := mocks.NewMockScheduler(ctrl)

    mockPodManager := mocks.NewMockPodManager(ctrl)
    orchestrator := NewDefaultPodOrchestrator(mockRepo, mockScheduler, mockPodManager)
	
    pod := &shared.Pod{ID: "pod1", Name: "test-pod"}

    // Success scenario
    mockScheduler.EXPECT().SchedulePod(pod).Return(nil)
    mockRepo.EXPECT().CreatePod(pod).Return(nil)
    mockPodManager.EXPECT().RunPod(gomock.Any()).Times(1)

	// Act
    err := orchestrator.OrchestratePodCreation(pod)

	// Assert
    assert.NoError(t, err)

    // Error in scheduling
    mockScheduler.EXPECT().SchedulePod(pod).Return(errors.New("scheduling failed"))
    
	// Act
	err = orchestrator.OrchestratePodCreation(pod)
    
	// Assert
	assert.Error(t, err)

    // Error in creating pod
    mockScheduler.EXPECT().SchedulePod(pod).Return(nil)
    mockRepo.EXPECT().CreatePod(pod).Return(errors.New("creation failed"))
    
	// Act
	err = orchestrator.OrchestratePodCreation(pod)
    
	// Assert
	assert.Error(t, err)
}

func TestOrchestratePodDeletion(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockPodRepository(ctrl)
    mockPodManager := mocks.NewMockPodManager(ctrl)
    orchestrator := NewDefaultPodOrchestrator(mockRepo, nil, mockPodManager)

    pod := &shared.Pod{ID: "pod1"}

    // Success scenario
    mockPodManager.EXPECT().StopPod(pod).Return(nil)
    mockRepo.EXPECT().DeletePod(pod.ID).Return(nil)

	// Act
    err := orchestrator.OrchestratePodDeletion(pod)

	// Assert
    assert.NoError(t, err)

    // Error in stopping pod
    mockPodManager.EXPECT().StopPod(pod).Return(errors.New("stopping failed"))
    
	// Act
	err = orchestrator.OrchestratePodDeletion(pod)

	// Assert
    assert.Error(t, err)

    // Error in deleting pod
    mockPodManager.EXPECT().StopPod(pod).Return(nil)
    mockRepo.EXPECT().DeletePod(pod.ID).Return(errors.New("deletion failed"))
    
	// Act
	err = orchestrator.OrchestratePodDeletion(pod)

	// Assert
    assert.Error(t, err)
}
