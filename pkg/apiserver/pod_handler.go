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

type PodHandler struct {
	Repo etcd.PodRepository
	orchestrator orchestrator.PodOrchestrator
}

func NewPodHandler(
	repo etcd.PodRepository,
	orchestrator orchestrator.PodOrchestrator,
) *PodHandler {
	return &PodHandler{Repo: repo, orchestrator: orchestrator}
}


func (h *PodHandler) createPodHandler(w http.ResponseWriter, r *http.Request) {
	var pod shared.Pod
	if err := json.NewDecoder(r.Body).Decode(&pod); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	err := h.orchestrator.OrchestratePodCreation(&pod)
	if err != nil {
		var dupErr *shared.ErrDuplicateResource
		if errors.As(err, &dupErr) {
			http.Error(w, dupErr.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pod)
}


func (h *PodHandler) listPodsHandler(w http.ResponseWriter, r *http.Request) {
	pods, err := h.Repo.ListPods()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pods)
}


func (h *PodHandler) deletePodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podID := vars["id"]

	if err := h.Repo.DeletePod(podID); err != nil {
		var notFoundErr *shared.ErrNotFound
		if errors.As(err, &notFoundErr) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
