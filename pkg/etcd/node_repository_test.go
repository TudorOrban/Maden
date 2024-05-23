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

func TestEtcdNodeRepositoryListNodes(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdNodeRepository(mockClient, mockTransactioner)

	mockClient.EXPECT().
		Get(gomock.Any(), nodesKey, gomock.Any()).
		Return(&clientv3.GetResponse{
			Kvs: []*mvccpb.KeyValue{
				{
					Key: []byte(nodesKey + "1"),
					Value: []byte(`{"id": "1", "name": "test-node"}`),
				},
			},
		}, nil).Times(1)

	// Act
	nodes, err := repo.ListNodes()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, "1", nodes[0].ID)
	assert.Equal(t, "test-node", nodes[0].Name)
}

func TestEtcdNodeRepositoryCreateNode(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl) 
    mockTransactioner := mocks.NewMockTransactioner(ctrl) 

    repo := NewEtcdNodeRepository(mockClient, mockTransactioner)

    node := &shared.Node{
        ID:   "1",
        Name: "test-node",
    }

    nodeData, _ := json.Marshal(node)
    key := nodesKey + node.ID

    mockTransactioner.EXPECT().
        PerformTransaction(gomock.Any(), key, string(nodeData), shared.NodeResource).
        Return(nil).Times(1)

    err := repo.CreateNode(node)
    assert.NoError(t, err)
}

func TestEtcdNodeRepositoryCreateNodeError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
    mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdNodeRepository(mockClient, mockTransactioner)

    node := &shared.Node{
        ID:   "1",
        Name: "test-node",
    }

    nodeData, _ := json.Marshal(node)
    key := nodesKey + node.ID

    mockTransactioner.EXPECT().
        PerformTransaction(gomock.Any(), key, string(nodeData), shared.NodeResource).
        Return(errors.New("transaction failed")).Times(1)

    err := repo.CreateNode(node)
    assert.Error(t, err)
    assert.Equal(t, "transaction failed", err.Error())
}

func TestEtcdNodeRepositoryUpdateNode(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdNodeRepository(mockClient, mockTransactioner)

    node := &shared.Node{ID: "1", Name: "updated-node"}

    mockClient.EXPECT().
        Put(gomock.Any(), nodesKey+node.ID, gomock.Any(), gomock.Any()).
        Return(&clientv3.PutResponse{PrevKv: &mvccpb.KeyValue{}}, nil).Times(1)

    // Act
    err := repo.UpdateNode(node)

    // Assert
    assert.NoError(t, err)
}

func TestEtcdNodeRepositoryUpdateNodeErrorOnPut(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdNodeRepository(mockClient, mockTransactioner)

    node := &shared.Node{ID: "1", Name: "updated-node"}

    mockClient.EXPECT().
        Put(gomock.Any(), nodesKey+node.ID, gomock.Any(), gomock.Any()).
        Return(nil, errors.New("etcd put error")).Times(1)

    err := repo.UpdateNode(node)
    assert.Error(t, err)
    assert.Equal(t, "etcd put error", err.Error())
}

func TestEtcdNodeRepositoryUpdateNodeNoPreviousKeyValue(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdNodeRepository(mockClient, mockTransactioner)

    node := &shared.Node{ID: "1", Name: "updated-node"}

    mockClient.EXPECT().
        Put(gomock.Any(), nodesKey+node.ID, gomock.Any(), gomock.Any()).
        Return(&clientv3.PutResponse{PrevKv: nil}, nil).Times(1)

    err := repo.UpdateNode(node)
    assert.Error(t, err)
    assert.IsType(t, &shared.ErrNotFound{}, err)
}

func TestEtcdNodeRepositoryDeleteNode(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdNodeRepository(mockClient, mockTransactioner)

    nodeID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), nodesKey+nodeID).
        Return(&clientv3.DeleteResponse{Deleted: 1}, nil).Times(1)

    // Act
    err := repo.DeleteNode(nodeID)

    // Assert
    assert.NoError(t, err)
}

func TestEtcdNodeRepositoryDeleteNodeErrorOnDelete(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdNodeRepository(mockClient, mockTransactioner)

    nodeID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), nodesKey+nodeID).
        Return(nil, errors.New("etcd delete error")).Times(1)

    err := repo.DeleteNode(nodeID)
    assert.Error(t, err)
    assert.Equal(t, "etcd delete error", err.Error())
}

func TestEtcdNodeRepositoryDeleteNodeNotFound(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockClient := mocks.NewMockEtcdClient(ctrl)
	mockTransactioner := mocks.NewMockTransactioner(ctrl)
    repo := NewEtcdNodeRepository(mockClient, mockTransactioner)

    nodeID := "1"

    mockClient.EXPECT().
        Delete(gomock.Any(), nodesKey+nodeID).
        Return(&clientv3.DeleteResponse{Deleted: 0}, nil).Times(1)

    err := repo.DeleteNode(nodeID)
    assert.Error(t, err)
    assert.IsType(t, &shared.ErrNotFound{}, err)
}
