package apiserver

import (
	"fmt"
	"log"
	"maden/pkg/controller"
	"net/http"

	"github.com/gorilla/mux"
)

func InitAPIServer() {
	registerRoutes()

	go controller.WatchDeployments()
	go controller.WatchServices()

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	
}

func registerRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/nodes", listNodesHandler).Methods("GET")
	r.HandleFunc("/nodes", createNodeHandler).Methods("POST")
	r.HandleFunc("/nodes/{id}", deleteNodeHandler).Methods("DELETE")
	r.HandleFunc("/pods", listPodsHandler).Methods("GET")
	r.HandleFunc("/pods", createPodHandler).Methods("POST")
	r.HandleFunc("/pods/{id}", deletePodHandler).Methods("DELETE")
	r.HandleFunc("/deployments", listDeploymentsHandler).Methods("GET")
	r.HandleFunc("/deployments/{name}", deleteDeploymentHandler).Methods("DELETE")
	r.HandleFunc("/services", listServicesHandler).Methods("GET")
	r.HandleFunc("/services/{name}", deleteServiceHandler).Methods("DELETE")
	r.HandleFunc("/manifests", handleMadenResources).Methods("POST")

	http.Handle("/", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Maden API Server")
}
