package controller

import (
	"maden/pkg/shared"
	"maden/pkg/mocks"
	
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleIncomingServiceCreateNew(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockServiceRepository(ctrl)
	mockDNSRepo := mocks.NewMockDNSRepository(ctrl)
	mockIPManager := mocks.NewMockIPManager(ctrl)
    controller := NewDefaultServiceController(mockRepo, mockDNSRepo, mockIPManager)

    serviceSpec := shared.ServiceSpec{Name: "test-service", Selector: map[string]string{"app": "myapp"}, Ports: []shared.ServicePort{{Port: 80, TargetPort: 8080}}}
    expectedService := transformToService(serviceSpec)

    mockRepo.EXPECT().GetServiceByName("test-service").Return(nil, &shared.ErrNotFound{})
    mockRepo.EXPECT().CreateService(gomock.Any()).Do(func(service *shared.Service) {
        assert.Equal(t, expectedService.Name, service.Name)
        assert.True(t, areMapsEqual(expectedService.Selector, service.Selector))
        assert.True(t, arePortsEqual(expectedService.Ports, service.Ports))
    }).Return(nil)
	mockDNSRepo.EXPECT().RegisterService("test-service", gomock.Any()).Return(nil)
	mockIPManager.EXPECT().AssignIP().Return(gomock.Any().String(), nil)

    // Act
    err := controller.HandleIncomingService(serviceSpec)

    // Assert
    assert.NoError(t, err)
}

func TestHandleIncomingServiceUpdateExisting(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockServiceRepository(ctrl)
	mockDNSRepo := mocks.NewMockDNSRepository(ctrl)
	mockIPManager := mocks.NewMockIPManager(ctrl)
    controller := NewDefaultServiceController(mockRepo, mockDNSRepo, mockIPManager)

    existingService := shared.Service{
        ID: "123",
        Name: "test-service",
        Selector: map[string]string{"app": "oldapp"},
        Ports: []shared.ServicePort{{Port: 80, TargetPort: 8080}},
    }
    serviceSpec := shared.ServiceSpec{
        Name: "test-service",
        Selector: map[string]string{"app": "newapp"},
        Ports: []shared.ServicePort{{Port: 80, TargetPort: 9090}},
    }

    mockRepo.EXPECT().GetServiceByName("test-service").Return(&existingService, nil)
    mockRepo.EXPECT().UpdateService(gomock.Any()).Do(func(service *shared.Service) {
        assert.Equal(t, "newapp", service.Selector["app"])
        assert.Equal(t, 9090, service.Ports[0].TargetPort)
    }).Return(nil)
	mockDNSRepo.EXPECT().RegisterService("test-service", gomock.Any()).Return(nil)

    // Act
    err := controller.HandleIncomingService(serviceSpec)

    // Assert
    assert.NoError(t, err)
}

func TestHandleIncomingServiceNoOperationNeeded(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockServiceRepository(ctrl)
	mockDNSRepo := mocks.NewMockDNSRepository(ctrl)
	mockIPManager := mocks.NewMockIPManager(ctrl)
    controller := NewDefaultServiceController(mockRepo, mockDNSRepo, mockIPManager)

    existingService := shared.Service{
        ID: "123",
        Name: "test-service",
        Selector: map[string]string{"app": "myapp"},
        Ports: []shared.ServicePort{{Port: 80, TargetPort: 8080}},
    }

    mockRepo.EXPECT().GetServiceByName("test-service").Return(&existingService, nil)

    // Act
    err := controller.HandleIncomingService(shared.ServiceSpec{Name: "test-service", Selector: map[string]string{"app": "myapp"}, Ports: []shared.ServicePort{{Port: 80, TargetPort: 8080}}})

    // Assert
    assert.NoError(t, err)
}
