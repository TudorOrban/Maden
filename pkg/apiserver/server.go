package apiserver

import (
	"fmt"
	"log"
	"maden/pkg/controller"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	PodHandler *PodHandler
	NodeHandler *NodeHandler
	DeploymentHandler *DeploymentHandler
	ServiceHandler *ServiceHandler
	ManifestHandler *ManifestHandler

	ChangeListener *controller.EtcdChangeListener
}

func NewServer(
	podHandler *PodHandler,
	nodeHandler *NodeHandler,
	deploymentHandler *DeploymentHandler, 
	serviceHandler *ServiceHandler,
	manifestHandler *ManifestHandler,
	changeListener *controller.EtcdChangeListener,
) *Server {
	s := &Server{
		router: mux.NewRouter(),
		PodHandler: podHandler,
		NodeHandler: nodeHandler,
		DeploymentHandler: deploymentHandler,
		ServiceHandler: serviceHandler,
		ManifestHandler: manifestHandler,
		ChangeListener: changeListener,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/", HomeHandler)
	s.router.HandleFunc("/pods", s.PodHandler.listPodsHandler).Methods("GET")
	s.router.HandleFunc("/pods", s.PodHandler.createPodHandler).Methods("POST")
	s.router.HandleFunc("/pods/{id}", s.PodHandler.deletePodHandler).Methods("DELETE")
	s.router.HandleFunc("/nodes", s.NodeHandler.listNodesHandler).Methods("GET")
	s.router.HandleFunc("/nodes", s.NodeHandler.createNodeHandler).Methods("POST")
	s.router.HandleFunc("/nodes/{id}", s.NodeHandler.deleteNodeHandler).Methods("DELETE")
	s.router.HandleFunc("/deployments", s.DeploymentHandler.listDeploymentsHandler).Methods("GET")
	s.router.HandleFunc("/deployments/{name}", s.DeploymentHandler.deleteDeploymentHandler).Methods("DELETE")
	s.router.HandleFunc("/services", s.ServiceHandler.listServicesHandler).Methods("GET")
	s.router.HandleFunc("/services/{name}", s.ServiceHandler.deleteServiceHandler).Methods("DELETE")
	s.router.HandleFunc("/manifests", s.ManifestHandler.handleMadenResources).Methods("POST")
}

func (s *Server) Start() {
	go s.ChangeListener.WatchDeployments()
	go s.ChangeListener.WatchServices()

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", s.router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Maden API Server")
}
