package apiserver

import (
	"maden/pkg/shared"
	"maden/pkg/etcd"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func createNodeHandler(w http.ResponseWriter, r *http.Request) {
	var node shared.Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := etcd.CreateNode(&node); err != nil {
		var dupErr *shared.ErrDuplicateResource
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


func listNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes, err := etcd.ListNodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}
