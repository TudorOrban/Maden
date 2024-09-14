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

type PersistentVolumeHandler struct {
	Repo       etcd.PersistentVolumeRepository
	Controller controller.PersistentVolumeController
}

func NewPersistentVolumeHandler(
	repo etcd.PersistentVolumeRepository,
	Controller controller.PersistentVolumeController,
) *PersistentVolumeHandler {
	return &PersistentVolumeHandler{Repo: repo, Controller: Controller}
}

func (h *PersistentVolumeHandler) listPersistentVolumesHandler(w http.ResponseWriter, r *http.Request) {
	persistentVolumes, err := h.Repo.ListPersistentVolumes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persistentVolumes)
}

func (h *PersistentVolumeHandler) deletePersistentVolumeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	persistentVolumeID := vars["iD"]

	if err := h.Repo.DeletePersistentVolume(persistentVolumeID); err != nil {
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