package etcd

import (
	"context"
	"encoding/json"
	"maden/pkg/shared"
	"time"

	"go.etcd.io/etcd/client/v3"
)

var nodesKey = "nodes/"

type EtcdNodeRepository struct {
	client EtcdClient
	transactioner Transactioner
}

func NewEtcdNodeRepository(
	client EtcdClient,
	transactioner Transactioner,
) NodeRepository {
	return &EtcdNodeRepository{client: client, transactioner: transactioner}
}


func (repo *EtcdNodeRepository) ListNodes() ([]shared.Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := repo.client.Get(ctx, nodesKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	nodes := make([]shared.Node, 0)
	for _, kv := range resp.Kvs {
		var node shared.Node
		if err := json.Unmarshal(kv.Value, &node); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (repo *EtcdNodeRepository) CreateNode(node *shared.Node) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	nodeData, err := json.Marshal(node)
	if err != nil {
		return err
	}

	key := nodesKey + node.ID

	return repo.transactioner.PerformTransaction(ctx, key, string(nodeData), shared.NodeResource)
}

func (repo *EtcdNodeRepository) UpdateNode(node *shared.Node) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    nodeData, err := json.Marshal(node)
    if err != nil {
        return err
    }

    key := nodesKey + node.ID

    resp, err := repo.client.Put(ctx, key, string(nodeData))
	if err != nil {
		return err
	}

	if resp.PrevKv == nil {
		return &shared.ErrNotFound{ID: node.ID, ResourceType: shared.NodeResource}
	}
    return err
}

func (repo *EtcdNodeRepository) DeleteNode(nodeID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := nodesKey + nodeID

	resp, err := repo.client.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return &shared.ErrNotFound{ID: nodeID, ResourceType: shared.NodeResource}
	}
	return nil
}