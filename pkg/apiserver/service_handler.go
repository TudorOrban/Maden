package apiserver

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"

	"encoding/json"
	"net/http"
	"errors"

	"github.com/gorilla/mux"
)

func listServicesHandler(w http.ResponseWriter, r *http.Request) {
	services, err := etcd.ListServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func deleteServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["name"]

	if err := etcd.DeleteService(serviceName); err != nil {
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