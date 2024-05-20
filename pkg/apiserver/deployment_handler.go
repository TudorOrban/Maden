package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maden/pkg/etcd"
	"maden/pkg/shared"
	"net/http"

	"gopkg.in/yaml.v3"
)

func listDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	deployments, err := etcd.ListDeployments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployments)
}


func createDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var deployment shared.Deployment
	if err := yaml.Unmarshal(body, &deployment); err != nil {
		http.Error(w, "Failed to parse YAML", http.StatusBadRequest)
		return
	}

	if err := etcd.CreateDeployment(&deployment); err != nil {
		var dupErr *shared.ErrDuplicateResource
		if errors.As(err, &dupErr) {
			http.Error(w, dupErr.Error(), http.StatusConflict)
		} else {
			http.Error(w, fmt.Sprintf("Error storing deployment data in etcd: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deployment)
}