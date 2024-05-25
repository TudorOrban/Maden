package controller

import (
	"maden/pkg/shared"
	"maden/pkg/mocks"
	
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleIncomingDeploymentCreateNew(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeploymentRepository(ctrl)
	controller := NewDefaultDeploymentController(mockRepo)

	deploymentSpec := shared.DeploymentSpec{Name: "test-deployment", Replicas: 3}
	expectedDeployment := transformToDeployment(deploymentSpec)

	mockRepo.EXPECT().GetDeploymentByName("test-deployment").Return(nil, &shared.ErrNotFound{})
	mockRepo.EXPECT().CreateDeployment(gomock.Any()).Do(func(deployment *shared.Deployment) {
		assert.Equal(t, expectedDeployment.Name, deployment.Name)
		assert.Equal(t, expectedDeployment.Replicas, deployment.Replicas)
	}).Return(nil)

	// Act
	err := controller.HandleIncomingDeployment(deploymentSpec)

	// Assert
	assert.NoError(t, err)
}

func TestHandleIncomingDeploymentUpdateExisting(t *testing.T) {
	// Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockDeploymentRepository(ctrl)
    controller := NewDefaultDeploymentController(mockRepo)

    existingDeployment := shared.Deployment{
        ID: "123",
        Name: "test-deployment",
        Replicas: 2,
        Selector: shared.LabelSelector{MatchLabels: map[string]string{"app": "old"}},
        Template: shared.PodTemplate{
            Metadata: shared.Metadata{Labels: map[string]string{"app": "old"}},
            Spec: shared.PodSpec{
                Containers: []shared.Container{{Image: "old-image", Ports: []shared.Port{{ContainerPort: 8080}}}},
            },
        },
    }
    deploymentSpec := shared.DeploymentSpec{
        Name: "test-deployment",
        Replicas: 3, // Different number of replicas to trigger an update
        Selector: shared.LabelSelector{MatchLabels: map[string]string{"app": "new"}},
        Template: shared.PodTemplate{
            Metadata: shared.Metadata{Labels: map[string]string{"app": "new"}},
            Spec: shared.PodSpec{
                Containers: []shared.Container{{Image: "new-image", Ports: []shared.Port{{ContainerPort: 9090}}}},
            },
        },
    }

    mockRepo.EXPECT().GetDeploymentByName("test-deployment").Return(&existingDeployment, nil)
    mockRepo.EXPECT().UpdateDeployment(gomock.Any()).Do(func(deployment *shared.Deployment) {
        assert.Equal(t, 3, deployment.Replicas)
        assert.Equal(t, "new-image", deployment.Template.Spec.Containers[0].Image)
    }).Return(nil)

	// Act
    err := controller.HandleIncomingDeployment(deploymentSpec)

    // Assert
    assert.NoError(t, err)
}


func TestHandleIncomingDeploymentNoOperationNeeded(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockDeploymentRepository(ctrl)
	controller := NewDefaultDeploymentController(mockRepo)

	existingDeployment := shared.Deployment{ID: "123", Name: "test-deployment", Replicas: 3}

	mockRepo.EXPECT().GetDeploymentByName("test-deployment").Return(&existingDeployment, nil)

	// Act
	err := controller.HandleIncomingDeployment(shared.DeploymentSpec{Name: "test-deployment", Replicas: 3})

	// Assert
	assert.NoError(t, err)
}
