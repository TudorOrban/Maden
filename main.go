package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"	
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/pods", createPodHandler).Methods("POST")
	r.HandleFunc("/pods", listPodsHandler).Methods("GET")
	r.HandleFunc("/pods/{id}", deletePodHandler).Methods("DELETE")
	http.Handle("/", r)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Maden API Server")
}

func createPodHandler(w http.ResponseWriter, r *http.Request) {
	var pod Pod
	if err := json.NewDecoder(r.Body).Decode(&pod); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

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

var podDB = []Pod{}
var mu sync.Mutex


type Pod struct {
	ID string `json:"id"`
	Name string `json:"name"`
}
