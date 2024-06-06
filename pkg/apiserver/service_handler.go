package apiserver

import (
	"maden/pkg/etcd"
	"maden/pkg/orchestrator"
	"maden/pkg/shared"

	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type ServiceHandler struct {
	Repo etcd.ServiceRepository
	SvcOrchestrator orchestrator.ServiceOrchestrator
}

func NewServiceHandler(
	repo etcd.ServiceRepository,
	svcOrchestrator orchestrator.ServiceOrchestrator,
	) *ServiceHandler {
	return &ServiceHandler{Repo: repo, SvcOrchestrator: svcOrchestrator}
}

func (h *ServiceHandler) listServicesHandler(w http.ResponseWriter, r *http.Request) {
	services, err := h.Repo.ListServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (h *ServiceHandler) deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["name"]

	if err := h.SvcOrchestrator.OrchestrateServiceDeletion(serviceName); err != nil {
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