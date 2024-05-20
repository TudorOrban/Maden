package main

import (
	// "context"
	"encoding/json"
	"net/http"
	// "time"
)

func registerNodeHandler(w http.ResponseWriter, r *http.Request) {
	var node Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	

	// Prevent duplicates
	for _, n := range nodeDB {
		if n.ID == node.ID || n.Name == node.Name {
			w.WriteHeader(http.StatusConflict)
			mu.Unlock()
			return
		}
	}

	nodeDB = append(nodeDB, node)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(node)
}

// func createNode(node *Node) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
// 	defer cancel()

// 	nodeData, err := json.Marshal(node)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = cli.Put(ctx, "nodes/" + node.ID, string(nodeData))
// 	return err
// }

func listNodesHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	json.NewEncoder(w).Encode(nodeDB)
	mu.Unlock()
}

