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

func TestEtcdPodRepositoryListPods(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdPodRepository(mockClient, mockTransactioner)

	mockClient.EXPECT().
		Get(gomock.Any(), podsKey, gomock.Any()).
		Return(&clientv3.GetResponse{
			Kvs: []*mvccpb.KeyValue{
				{
					Key: []byte(podsKey + "1"),
					Value: []byte(`{"id": "1", "name": "test-pod", "deploymentID": "1"}`),
				},
			},
		}, nil).Times(1)

	// Act
	pods, err := repo.ListPods()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, pods, 1)
	assert.Equal(t, "1", pods[0].ID)
	assert.Equal(t, "test-pod", pods[0].Name)
	assert.Equal(t, "1", pods[0].DeploymentID)
}

func TestEtcdPodRepositoryGetPodsByDeploymentID(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdPodRepository(mockClient, mockTransactioner)

    deploymentID := "1"

    mockClient.EXPECT().
        Get(gomock.Any(), podsKey, gomock.Any()).
        Return(&clientv3.GetResponse{
            Kvs: []*mvccpb.KeyValue{
                {
                    Key:   []byte(podsKey + "1"),
                    Value: []byte(`{"id": "1", "name": "test-pod", "deploymentID": "1"}`),
                },
            },
        }, nil).Times(1)

    // Act
    pods, err := repo.GetPodsByDeploymentID(deploymentID)

    // Assert
    assert.NoError(t, err)
    assert.Len(t, pods, 1) 
    assert.Equal(t, deploymentID, pods[0].DeploymentID)
}

func TestEtcdPodRepositoryCreatePod(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl) 
    mockTransactioner := mocks.NewMockTransactioner(ctrl) 

    repo := NewEtcdPodRepository(mockClient, mockTransactioner)

    pod := &shared.Pod{
        ID:   "1",
        Name: "test-pod",
    }

    podData, _ := json.Marshal(pod)
    key := podsKey + pod.ID

    mockTransactioner.EXPECT().
        PerformTransaction(gomock.Any(), key, string(podData), shared.PodResource).
        Return(nil).Times(1)

    err := repo.CreatePod(pod)
    assert.NoError(t, err)
}

func TestEtcdPodRepositoryCreatePodError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
    mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdPodRepository(mockClient, mockTransactioner)

    pod := &shared.Pod{
        ID:   "1",
        Name: "test-pod",
    }

    podData, _ := json.Marshal(pod)
    key := podsKey + pod.ID

    mockTransactioner.EXPECT().
        PerformTransaction(gomock.Any(), key, string(podData), shared.PodResource).
        Return(errors.New("transaction failed")).Times(1)

    err := repo.CreatePod(pod)
    assert.Error(t, err)
    assert.Equal(t, "transaction failed", err.Error())
}

func TestEtcdPodRepositoryUpdatePod(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdPodRepository(mockClient, mockTransactioner)

    pod := &shared.Pod{ID: "1", Name: "updated-pod"}

    mockClient.EXPECT().
        Put(gomock.Any(), podsKey+pod.ID, gomock.Any(), gomock.Any()).
        Return(&clientv3.PutResponse{PrevKv: &mvccpb.KeyValue{}}, nil).Times(1)

    // Act
    err := repo.UpdatePod(pod)

    // Assert
    assert.NoError(t, err)
}

func TestEtcdPodRepositoryUpdatePodErrorOnPut(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdPodRepository(mockClient, mockTransactioner)

    pod := &shared.Pod{ID: "1", Name: "updated-pod"}

    mockClient.EXPECT().
        Put(gomock.Any(), podsKey+pod.ID, gomock.Any(), gomock.Any()).
        Return(nil, errors.New("etcd put error")).Times(1)

    err := repo.UpdatePod(pod)
    assert.Error(t, err)
    assert.Equal(t, "etcd put error", err.Error())
}

func TestEtcdPodRepositoryUpdatePodNoPreviousKeyValue(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdPodRepository(mockClient, mockTransactioner)

    pod := &shared.Pod{ID: "1", Name: "updated-pod"}

    mockClient.EXPECT().
        Put(gomock.Any(), podsKey+pod.ID, gomock.Any(), gomock.Any()).
        Return(&clientv3.PutResponse{PrevKv: nil}, nil).Times(1)

    err := repo.UpdatePod(pod)
    assert.Error(t, err)
    assert.IsType(t, &shared.ErrNotFound{}, err)
}

func TestEtcdPodRepositoryDeletePod(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdPodRepository(mockClient, mockTransactioner)

    podID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), podsKey+podID).
        Return(&clientv3.DeleteResponse{Deleted: 1}, nil).Times(1)

    // Act
    err := repo.DeletePod(podID)

    // Assert
    assert.NoError(t, err)
}

func TestEtcdPodRepositoryDeletePodErrorOnDelete(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdPodRepository(mockClient, mockTransactioner)

    podID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), podsKey+podID).
        Return(nil, errors.New("etcd delete error")).Times(1)

    err := repo.DeletePod(podID)
    assert.Error(t, err)
    assert.Equal(t, "etcd delete error", err.Error())
}

func TestEtcdPodRepositoryDeletePodNotFound(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdPodRepository(mockClient, mockTransactioner)

    podID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), podsKey+podID).
        Return(&clientv3.DeleteResponse{Deleted: 0}, nil).Times(1)

    err := repo.DeletePod(podID)
    assert.Error(t, err)
    assert.IsType(t, &shared.ErrNotFound{}, err)
}
