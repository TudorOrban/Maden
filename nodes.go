package main

import (
	"encoding/json"
	"net/http"
)

func registerNodeHandler(w http.ResponseWriter, r *http.Request) {
	var node Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()

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

func listNodesHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	json.NewEncoder(w).Encode(nodeDB)
	mu.Unlock()
}

