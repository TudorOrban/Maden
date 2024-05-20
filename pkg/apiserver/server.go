package apiserver

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func InitAPIServer() {
	registerRoutes()

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func registerRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/nodes", listNodesHandler).Methods("GET")
	r.HandleFunc("/nodes", createNodeHandler).Methods("POST")
	r.HandleFunc("/pods", listPodsHandler).Methods("GET")
	r.HandleFunc("/pods", createPodHandler).Methods("POST")
	r.HandleFunc("/pods/{id}", deletePodHandler).Methods("DELETE")
	r.HandleFunc("/deployments", listDeploymentsHandler).Methods("GET")
	r.HandleFunc("/deployments", createDeploymentHandler).Methods("POST")
	r.HandleFunc("/deployments/{id}", deleteDeploymentHandler).Methods("DELETE")

	http.Handle("/", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Maden API Server")
}
