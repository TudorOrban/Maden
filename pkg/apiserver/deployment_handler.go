package apiserver

import (
	"maden/pkg/controller"
	"maden/pkg/etcd"
	"maden/pkg/shared"

	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type DeploymentHandler struct {
	Repo etcd.DeploymentRepository
	UpdateController controller.DeploymentUpdaterController
}

func NewDeploymentHandler(
	repo etcd.DeploymentRepository,
	updateController controller.DeploymentUpdaterController,
	) *DeploymentHandler {
	return &DeploymentHandler{Repo: repo, UpdateController: updateController}
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

func (h *DeploymentHandler) deleteDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deploymentName := vars["name"]

	if err := h.Repo.DeleteDeployment(deploymentName); err != nil {
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

func (h *DeploymentHandler) rolloutRestartDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deploymentName := vars["name"]

	deployment, err := h.Repo.GetDeploymentByName(deploymentName)

	if err != nil {
		var errNotFound *shared.ErrNotFound
		if errors.As(err, &errNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = h.UpdateController.HandleDeploymentRolloutRestart(deployment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DeploymentHandler) scaleDeploymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deploymentName := vars["name"]

	deployment, err := h.Repo.GetDeploymentByName(deploymentName)
	if err != nil {
		var errNotFound *shared.ErrNotFound
		if errors.As(err, &errNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var scaleRequest shared.ScaleRequest
	if err := json.NewDecoder(r.Body).Decode(&scaleRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deployment.Replicas = scaleRequest.Replicas

	err = h.Repo.UpdateDeployment(deployment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}