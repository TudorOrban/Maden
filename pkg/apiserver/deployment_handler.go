package apiserver

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"

	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
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

	
    fmt.Println("Received YAML:", string(body))

	var deploymentSpec shared.DeploymentSpec
    if err := yaml.Unmarshal(body, &deploymentSpec); err != nil {
        http.Error(w, "Failed to parse YAML: "+err.Error(), http.StatusBadRequest)
        return
    }

	if err := etcd.CreateDeployment(&deploymentSpec.Deployment); err != nil {
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
	json.NewEncoder(w).Encode(deploymentSpec.Deployment)
}

func deleteDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deploymentID := vars["id"]

	if err := etcd.DeleteDeployment(deploymentID); err != nil {
		if err == shared.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
