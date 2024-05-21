package etcd

import (
	"context"
	"encoding/json"
	"maden/pkg/shared"
	"time"

	"go.etcd.io/etcd/client/v3"
)

var nodesKey = "nodes/"


func ListNodes() ([]shared.Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, nodesKey, clientv3.WithPrefix())
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

func CreateNode(node *shared.Node) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	nodeData, err := json.Marshal(node)
	if err != nil {
		return err
	}

	key := nodesKey + node.ID

	// Start transaction to prevent duplicates
	txnResp, err := cli.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, string(nodeData))).
		Else(clientv3.OpGet(key)).
		Commit()

	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return &shared.ErrDuplicateResource{ID: node.ID, ResourceType: shared.NodeResource}
	}

	return nil
}

func UpdateNode(node *shared.Node) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    nodeData, err := json.Marshal(node)
    if err != nil {
        return err
    }

    key := nodesKey + node.ID
    _, err = cli.Put(ctx, key, string(nodeData))
    return err
}

func DeleteNode(nodeID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := nodesKey + nodeID

	resp, err := cli.Delete(ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return &shared.ErrNotFound{ID: nodeID, ResourceType: shared.NodeResource}
	}
	return nil
}