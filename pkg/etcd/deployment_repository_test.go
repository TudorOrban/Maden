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

func TestEtcdDeploymentRepositoryListDeployments(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdDeploymentRepository(mockClient, mockTransactioner)

	mockClient.EXPECT().
		Get(gomock.Any(), deploymentsKey, gomock.Any()).
		Return(&clientv3.GetResponse{
			Kvs: []*mvccpb.KeyValue{
				{
					Key: []byte(deploymentsKey + "1"),
					Value: []byte(`{"id": "1", "name": "test-deployment"}`),
				},
			},
		}, nil).Times(1)

	// Act
	deployments, err := repo.ListDeployments()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, deployments, 1)
	assert.Equal(t, "1", deployments[0].ID)
	assert.Equal(t, "test-deployment", deployments[0].Name)
}

func TestEtcdDeploymentRepositoryCreateDeployment(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl) 
    mockTransactioner := mocks.NewMockTransactioner(ctrl) 

    repo := NewEtcdDeploymentRepository(mockClient, mockTransactioner)

    deployment := &shared.Deployment{
        ID:   "1",
        Name: "test-deployment",
    }

    deploymentData, _ := json.Marshal(deployment)
    key := deploymentsKey + deployment.Name

    mockTransactioner.EXPECT().
        PerformTransaction(gomock.Any(), key, string(deploymentData), shared.DeploymentResource).
        Return(nil).Times(1)

    err := repo.CreateDeployment(deployment)
    assert.NoError(t, err)
}

func TestEtcdDeploymentRepositoryCreateDeploymentError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
    mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdDeploymentRepository(mockClient, mockTransactioner)

    deployment := &shared.Deployment{
        ID:   "1",
        Name: "test-deployment",
    }

    deploymentData, _ := json.Marshal(deployment)
    key := deploymentsKey + deployment.Name

    mockTransactioner.EXPECT().
        PerformTransaction(gomock.Any(), key, string(deploymentData), shared.DeploymentResource).
        Return(errors.New("transaction failed")).Times(1)

    err := repo.CreateDeployment(deployment)
    assert.Error(t, err)
    assert.Equal(t, "transaction failed", err.Error())
}

func TestEtcdDeploymentRepositoryUpdateDeployment(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdDeploymentRepository(mockClient, mockTransactioner)

    deployment := &shared.Deployment{ID: "1", Name: "updated-deployment"}

    mockClient.EXPECT().
        Put(gomock.Any(), deploymentsKey + deployment.Name, gomock.Any()).
        Return(&clientv3.PutResponse{PrevKv: &mvccpb.KeyValue{}}, nil).Times(1)

    // Act
    err := repo.UpdateDeployment(deployment)

    // Assert
    assert.NoError(t, err)
}

func TestEtcdDeploymentRepositoryUpdateDeploymentErrorOnPut(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdDeploymentRepository(mockClient, mockTransactioner)

    deployment := &shared.Deployment{ID: "1", Name: "updated-deployment"}

    mockClient.EXPECT().
        Put(gomock.Any(), deploymentsKey + deployment.Name, gomock.Any()).
        Return(nil, errors.New("etcd put error")).Times(1)

    err := repo.UpdateDeployment(deployment)
    assert.Error(t, err)
    assert.Equal(t, "etcd put error", err.Error())
}

func TestEtcdDeploymentRepositoryUpdateDeploymentNoPreviousKeyValue(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdDeploymentRepository(mockClient, mockTransactioner)

    deployment := &shared.Deployment{ID: "1", Name: "updated-deployment"}

    mockClient.EXPECT().
        Put(gomock.Any(), deploymentsKey+deployment.Name, gomock.Any()).
        Return(&clientv3.PutResponse{PrevKv: nil}, nil).Times(1)

    err := repo.UpdateDeployment(deployment)
    assert.Error(t, err)
    assert.IsType(t, &shared.ErrNotFound{}, err)
}

func TestEtcdDeploymentRepositoryDeleteDeployment(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdDeploymentRepository(mockClient, mockTransactioner)

    deploymentID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), deploymentsKey+deploymentID).
        Return(&clientv3.DeleteResponse{Deleted: 1}, nil).Times(1)

    // Act
    err := repo.DeleteDeployment(deploymentID)

    // Assert
    assert.NoError(t, err)
}

func TestEtcdDeploymentRepositoryDeleteDeploymentErrorOnDelete(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdDeploymentRepository(mockClient, mockTransactioner)

    deploymentID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), deploymentsKey+deploymentID).
        Return(nil, errors.New("etcd delete error")).Times(1)

    err := repo.DeleteDeployment(deploymentID)
    assert.Error(t, err)
    assert.Equal(t, "etcd delete error", err.Error())
}

func TestEtcdDeploymentRepositoryDeleteDeploymentNotFound(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdDeploymentRepository(mockClient, mockTransactioner)

    deploymentID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), deploymentsKey+deploymentID).
        Return(&clientv3.DeleteResponse{Deleted: 0}, nil).Times(1)

    err := repo.DeleteDeployment(deploymentID)
    assert.Error(t, err)
    assert.IsType(t, &shared.ErrNotFound{}, err)
}
