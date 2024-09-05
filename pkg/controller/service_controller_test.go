package controller

import (
	"maden/pkg/mocks"
	"maden/pkg/shared"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)


func TestHandleIncomingService(t *testing.T) {
 // Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	mockOrchestrator := mocks.NewMockServiceOrchestrator(ctrl)
	serviceController := NewDefaultServiceController(mockRepo, mockOrchestrator)
	
	serviceSpec := shared.ServiceSpec{
		Name: "test-service",
		Selector: map[string]string{"app": "myapp"},
		Ports: []shared.ServicePort{{Port: 80, TargetPort: 8080}},
	}

	t.Run("Service Creation", func(t *testing.T) {
		notFoundErr := &shared.ErrNotFound{Name: serviceSpec.Name, ResourceType: shared.ServiceResource}
        mockRepo.EXPECT().GetServiceByName(gomock.Eq(serviceSpec.Name)).Return(nil, notFoundErr)
        mockOrchestrator.EXPECT().OrchestrateServiceCreation(serviceSpec).Return(nil)

// Act
		err := serviceController.HandleIncomingService(serviceSpec)
		assert.NoError(t, err)
	})

	t.Run("Service Update", func(t *testing.T) {
		existingService := &shared.Service{
			Name: serviceSpec.Name,
			Selector: map[string]string{"app": "old"},
			Ports: []shared.ServicePort{{Port: 80, TargetPort: 8081}},
		}
	
// Assert	mockRepo.EXPECT().GetServiceByName(serviceSpec.Name).Return(existingService, nil)
		mockOrchestrator.EXPECT().OrchestrateServiceUpdate(*existingService, serviceSpec).Return(nil)

		err := serviceController.HandleIncomingService(serviceSpec)
		assert.NoError(t, err)
	})

	t.Run("No Update Required", func(t *testing.T) {
		existingService := &shared.Service{
			Name: serviceSpec.Name,
			Selector: serviceSpec.Selector,
			Ports: serviceSpec.Ports,
		}
		mockRepo.EXPECT().GetServiceByName(serviceSpec.Name).Return(existingService, nil)

		err := serviceController.HandleIncomingService(serviceSpec)
		assert.NoError(t, err)
	})
}
