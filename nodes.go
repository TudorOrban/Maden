package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func createNodeHandler(w http.ResponseWriter, r *http.Request) {
	var node Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := createNode(&node); err != nil {
		var dupErr *ErrDuplicateResource
		if errors.As(err, &dupErr) {
			http.Error(w, dupErr.Error(), http.StatusConflict)
		} else {
			http.Error(w, fmt.Sprintf("Error storing node data in etcd: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(node)
}

func createNode(node *Node) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	nodeData, err := json.Marshal(node)
	if err != nil {
		return err
	}

	key := "nodes/" + node.ID

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
		return &ErrDuplicateResource{ID: node.ID, ResourceType: NodeResource}
	}

	return nil
}

func listNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes, err := listNodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

func listNodes() ([]Node, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, "nodes/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, 0)
	for _, kv := range resp.Kvs {
		var node Node
		if err := json.Unmarshal(kv.Value, &node); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func updateNode(node *Node) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    nodeData, err := json.Marshal(node)
    if err != nil {
        return err
    }

    key := "nodes/" + node.ID
    _, err = cli.Put(ctx, key, string(nodeData))
    return err
}
