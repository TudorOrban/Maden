package etcd

import (
	"encoding/json"
	"errors"
	"maden/pkg/mocks"
	"maden/pkg/shared"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdServiceRepositoryListServices(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdServiceRepository(mockClient, mockTransactioner)

	mockClient.EXPECT().
		Get(gomock.Any(), servicesKey, gomock.Any()).
		Return(&clientv3.GetResponse{
			Kvs: []*mvccpb.KeyValue{
				{
					Key: []byte(servicesKey + "1"),
					Value: []byte(`{"id": "1", "name": "test-service"}`),
				},
			},
		}, nil).Times(1)

	// Act
	services, err := repo.ListServices()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, services, 1)
	assert.Equal(t, "1", services[0].ID)
	assert.Equal(t, "test-service", services[0].Name)
}

func TestEtcdServiceRepositoryCreateService(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl) 
    mockTransactioner := mocks.NewMockTransactioner(ctrl) 

    repo := NewEtcdServiceRepository(mockClient, mockTransactioner)

    service := &shared.Service{
        ID:   "1",
        Name: "test-service",
    }

    serviceData, _ := json.Marshal(service)
    key := servicesKey + service.Name

    mockTransactioner.EXPECT().
        PerformTransaction(gomock.Any(), key, string(serviceData), shared.ServiceResource).
        Return(nil).Times(1)

    err := repo.CreateService(service)
    assert.NoError(t, err)
}

func TestEtcdServiceRepositoryCreateServiceError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
    mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdServiceRepository(mockClient, mockTransactioner)

    service := &shared.Service{
        ID:   "1",
        Name: "test-service",
    }

    serviceData, _ := json.Marshal(service)
    key := servicesKey + service.Name

    mockTransactioner.EXPECT().
        PerformTransaction(gomock.Any(), key, string(serviceData), shared.ServiceResource).
        Return(errors.New("transaction failed")).Times(1)

    err := repo.CreateService(service)
    assert.Error(t, err)
    assert.Equal(t, "transaction failed", err.Error())
}

func TestEtcdServiceRepositoryUpdateService(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdServiceRepository(mockClient, mockTransactioner)

    service := &shared.Service{ID: "1", Name: "updated-service"}

    mockClient.EXPECT().
        Put(gomock.Any(), servicesKey + service.Name, gomock.Any(), gomock.Any()).
        Return(&clientv3.PutResponse{PrevKv: &mvccpb.KeyValue{}}, nil).Times(1)

    // Act
    err := repo.UpdateService(service)

    // Assert
    assert.NoError(t, err)
}

func TestEtcdServiceRepositoryUpdateServiceErrorOnPut(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdServiceRepository(mockClient, mockTransactioner)

    service := &shared.Service{ID: "1", Name: "updated-service"}

    mockClient.EXPECT().
        Put(gomock.Any(), servicesKey + service.Name, gomock.Any(), gomock.Any()).
        Return(nil, errors.New("etcd put error")).Times(1)

    err := repo.UpdateService(service)
    assert.Error(t, err)
    assert.Equal(t, "etcd put error", err.Error())
}

func TestEtcdServiceRepositoryUpdateServiceNoPreviousKeyValue(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdServiceRepository(mockClient, mockTransactioner)

    service := &shared.Service{ID: "1", Name: "updated-service"}

    mockClient.EXPECT().
        Put(gomock.Any(), servicesKey+service.Name, gomock.Any(), gomock.Any()).
        Return(&clientv3.PutResponse{PrevKv: nil}, nil).Times(1)

    err := repo.UpdateService(service)
    assert.Error(t, err)
    assert.IsType(t, &shared.ErrNotFound{}, err)
}

func TestEtcdServiceRepositoryDeleteService(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdServiceRepository(mockClient, mockTransactioner)

    serviceID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), servicesKey+serviceID).
        Return(&clientv3.DeleteResponse{Deleted: 1}, nil).Times(1)

    // Act
    err := repo.DeleteService(serviceID)

    // Assert
    assert.NoError(t, err)
}

func TestEtcdServiceRepositoryDeleteServiceErrorOnDelete(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdServiceRepository(mockClient, mockTransactioner)

    serviceID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), servicesKey+serviceID).
        Return(nil, errors.New("etcd delete error")).Times(1)

    err := repo.DeleteService(serviceID)
    assert.Error(t, err)
    assert.Equal(t, "etcd delete error", err.Error())
}

func TestEtcdServiceRepositoryDeleteServiceNotFound(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdServiceRepository(mockClient, mockTransactioner)

    serviceID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), servicesKey+serviceID).
        Return(&clientv3.DeleteResponse{Deleted: 0}, nil).Times(1)

    err := repo.DeleteService(serviceID)
    assert.Error(t, err)
    assert.IsType(t, &shared.ErrNotFound{}, err)
}
