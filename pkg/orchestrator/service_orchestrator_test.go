package orchestrator

import (
	"maden/pkg/shared"
	"maden/pkg/mocks"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDefaultServiceOrchestratorOrchestrateServiceCreation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	mockDNSRepo := mocks.NewMockDNSRepository(ctrl)
	mockIPManager := mocks.NewMockIPManager(ctrl)

	orchestrator := NewDefaultServiceOrchestrator(mockRepo, mockDNSRepo, mockIPManager)

	serviceSpec := shared.ServiceSpec{
		Name: "test-service",
		Selector: map[string]string{"app": "myapp"},
		Ports: []shared.ServicePort{{Port: 80, TargetPort: 8080}},
	}

	// Setting up the test scenario
	mockIPManager.EXPECT().AssignIP().Return("192.168.1.100", nil)
	mockRepo.EXPECT().CreateService(gomock.Any()).Return(nil)
	mockDNSRepo.EXPECT().RegisterService("test-service", "192.168.1.100").Return(nil)

	err := orchestrator.OrchestrateServiceCreation(serviceSpec)
	assert.NoError(t, err)
}

func TestDefaultServiceOrchestratorOrchestrateServiceDeletion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	mockDNSRepo := mocks.NewMockDNSRepository(ctrl)
	mockIPManager := mocks.NewMockIPManager(ctrl)

	orchestrator := NewDefaultServiceOrchestrator(mockRepo, mockDNSRepo, mockIPManager)

	serviceName := "test-service"
	service := shared.Service{
		ID: "1",
		Name: serviceName,
		IP: "192.168.1.100",
	}

	// Setting up the test scenario
	mockRepo.EXPECT().GetServiceByName(serviceName).Return(&service, nil)
	mockDNSRepo.EXPECT().DeregisterService(serviceName).Return(nil)
	mockIPManager.EXPECT().ReleaseIP("192.168.1.100").Return(nil)
	mockRepo.EXPECT().DeleteService(serviceName).Return(nil)

	err := orchestrator.OrchestrateServiceDeletion(serviceName)
	assert.NoError(t, err)
}

func TestDefaultServiceOrchestratorOrchestrateServiceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockServiceRepository(ctrl)
	mockDNSRepo := mocks.NewMockDNSRepository(ctrl)
	mockIPManager := mocks.NewMockIPManager(ctrl) // Might not be needed for the update scenario

	orchestrator := NewDefaultServiceOrchestrator(mockRepo, mockDNSRepo, mockIPManager)

	existingService := shared.Service{
		ID: "1",
		Name: "test-service",
		Selector: map[string]string{"app": "old"},
		Ports: []shared.ServicePort{{Port: 80, TargetPort: 8080}},
		IP: "192.168.1.100",
	}
	serviceSpec := shared.ServiceSpec{
		Name: "test-service",
		Selector: map[string]string{"app": "new"},
		Ports: []shared.ServicePort{{Port: 80, TargetPort: 8080}},
	}

	// Setting up the test scenario
	mockRepo.EXPECT().UpdateService(gomock.Any()).Return(nil)
	mockDNSRepo.EXPECT().RegisterService("test-service", "192.168.1.100").Return(nil)

	err := orchestrator.OrchestrateServiceUpdate(existingService, serviceSpec)
	assert.NoError(t, err)
}
