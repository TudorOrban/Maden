package apiserver

import (
	"maden/pkg/etcd"
	"maden/pkg/madelet"
	"maden/pkg/scheduler"
	"maden/pkg/shared"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func createPodHandler(w http.ResponseWriter, r *http.Request) {
	var pod shared.Pod
	if err := json.NewDecoder(r.Body).Decode(&pod); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := scheduler.SchedulePod(&pod)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error scheduling pod: %v", err), http.StatusInternalServerError)
	}

	if err := etcd.CreatePod(&pod); err != nil {
		var dupErr *shared.ErrDuplicateResource
		if errors.As(err, &dupErr) {
			http.Error(w, dupErr.Error(), http.StatusConflict)
		} else {
			http.Error(w, fmt.Sprintf("Error storing pod data in etcd: %v", err), http.StatusInternalServerError)
		}
		return
	}

	dockerRuntime := &madelet.DockerRuntime{}
	podLifecycleManager := madelet.PodLifecycleManager{
		Runtime: dockerRuntime,
	}
	go podLifecycleManager.RunPod(&pod)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pod)
}


func listPodsHandler(w http.ResponseWriter, r *http.Request) {
	pods, err := etcd.ListPods()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pods)
}


func deletePodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podID := vars["id"]

	if err := etcd.DeletePod(podID); err != nil {
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
