package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func createPodHandler(w http.ResponseWriter, r *http.Request) {
	var pod Pod
	if err := json.NewDecoder(r.Body).Decode(&pod); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	schedulePod(&pod)

	mu.Lock()
	podDB = append(podDB, pod)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pod)
}

func listPodsHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	json.NewEncoder(w).Encode(podDB)
	mu.Unlock()
}

func deletePodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	for index, pod := range podDB {
		if pod.ID == id {
			podDB = append(podDB[:index], podDB[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			mu.Unlock()
			return
		}
	}
	mu.Unlock()

	w.WriteHeader(http.StatusNotFound)
}