package apiserver

import (
	"io"
	"log"
	"maden/pkg/etcd"
	"maden/pkg/orchestrator"
	"maden/pkg/shared"
	"os"

	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type PodHandler struct {
	Repo etcd.PodRepository
	Orchestrator orchestrator.PodOrchestrator
}

func NewPodHandler(
	repo etcd.PodRepository,
	orchestrator orchestrator.PodOrchestrator,
) *PodHandler {
	return &PodHandler{Repo: repo, Orchestrator: orchestrator}
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

func (h *PodHandler) createPodHandler(w http.ResponseWriter, r *http.Request) {
	var pod shared.Pod
	if err := json.NewDecoder(r.Body).Decode(&pod); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	err := h.Orchestrator.OrchestratePodCreation(&pod)
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

func (h *PodHandler) deletePodHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podID := vars["id"]

	pod, err := h.Repo.GetPodByID(podID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if err := h.Orchestrator.OrchestratePodDeletion(pod); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (h *PodHandler) getPodLogsHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    podID := vars["id"]
    containerID := r.URL.Query().Get("containerID")
    follow := r.URL.Query().Get("follow") == "true"

    log.Printf("Request to stream logs: PodID=%s, ContainerID=%s, Follow=%t", podID, containerID, follow)

    logsReader, err := h.Orchestrator.GetPodLogs(podID, containerID, follow)
    if err != nil {
        log.Printf("Failed to retrieve logs: %v", err)
        http.Error(w, "Failed to get logs", http.StatusInternalServerError)
        return
	}
	defer logsReader.Close()

	if follow {
		_, err = io.Copy(os.Stdout, logsReader)
		if err != nil {
			log.Printf("Failed to stream logs for container %s: %v", containerID, err)
			return
		}
	} else {
		logContents, err := io.ReadAll(logsReader)
		if err != nil {
			log.Printf("Failed to read logs for container %s: %v", containerID, err)
			return
		}

		log.Printf("Logs for container %s: %s", containerID, logContents)
	}

    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Connection", "keep-alive")
    flusher, ok := w.(http.Flusher)
    if ! ok {
        http.Error(w, "Streaming not supported", http.StatusInternalServerError)
        return
    }

    ctx := r.Context()
    w.WriteHeader(http.StatusOK)
    flusher.Flush()

    buf := make([]byte, 1024)
    for {
        select {
        case <-ctx.Done():
            log.Printf("HTTP context was canceled, reason: %v", ctx.Err())
            return
        default:
            n, readErr := logsReader.Read(buf)
            if n > 0 {
                _, writeErr := w.Write(buf[:n])
                if writeErr != nil {
                    log.Printf("Failed to write logs: %v", writeErr)
                    return
                }
                flusher.Flush()
            }
            if readErr != nil {
                if readErr == io.EOF {
                    log.Println("Reached EOF for logs stream")
                    return
                }
                log.Printf("Error reading logs: %v", readErr)
                return
            }
        }
    }
}
