package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.etcd.io/etcd/client/v3"
)

func createPodHandler(w http.ResponseWriter, r *http.Request) {
	var pod Pod
	if err := json.NewDecoder(r.Body).Decode(&pod); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	// schedulePod(&pod)

	if err := storePod(&pod); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error storing pod data in etcd: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pod)
}

func storePod(pod *Pod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	podData, err := json.Marshal(pod)
	if err != nil {
		return err
	}

	_, err = cli.Put(ctx, "pods/" + pod.ID, string(podData))
	return err
}

func listPodsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, "pods/", clientv3.WithPrefix())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pods := make([]Pod, 0)
	for _, kv := range resp.Kvs {
		var pod Pod
		if err := json.Unmarshal(kv.Value, &pod); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pods = append(pods, pod)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pods); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deletePodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	deletePod(w, id)
}

func deletePod(w http.ResponseWriter, podID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	key := "pods/" + podID

	resp, err := cli.Delete(ctx, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.Deleted == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}