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

type PersistentVolumeClaimHandler struct {
	Repo       etcd.PersistentVolumeClaimRepository
	Controller controller.PersistentVolumeClaimController
}

func NewPersistentVolumeClaimHandler(
	repo etcd.PersistentVolumeClaimRepository,
	Controller controller.PersistentVolumeClaimController,
) *PersistentVolumeClaimHandler {
	return &PersistentVolumeClaimHandler{Repo: repo, Controller: Controller}
}

func (h *PersistentVolumeClaimHandler) listPersistentVolumeClaimsHandler(w http.ResponseWriter, r *http.Request) {
	persistentVolumeClaims, err := h.Repo.ListPersistentVolumeClaims()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persistentVolumeClaims)
}

func (h *PersistentVolumeClaimHandler) deletePersistentVolumeClaimHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	persistentVolumeClaimID := vars["iD"]

	if err := h.Repo.DeletePersistentVolumeClaim(persistentVolumeClaimID); err != nil {
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
