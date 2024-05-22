package apiserver

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"

	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type DeploymentHandler struct {
	Repo etcd.DeploymentRepository
}

func NewDeploymentHandler(repo etcd.DeploymentRepository) *DeploymentHandler {
	return &DeploymentHandler{Repo: repo}
}

func (h *DeploymentHandler) listDeploymentsHandler(w http.ResponseWriter, r *http.Request) {
	deployments, err := h.Repo.ListDeployments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployments)
}

func deleteDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deploymentName := vars["name"]

	if err := etcd.DeleteDeployment(deploymentName); err != nil {
		var errNotFound *shared.ErrNotFound
		if errors.As(err, &errNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

