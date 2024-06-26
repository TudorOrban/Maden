package apiserver

import (
	"fmt"
	"log"
	"maden/pkg/controller"
	"net/http"
	"time"

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
	s.router.HandleFunc("/pods/{id}/logs", s.PodHandler.getPodLogsHandler).Methods("GET")
	s.router.HandleFunc("/pods/{id}/exec", s.PodHandler.execWebSocketHandler).Methods("GET")
	s.router.HandleFunc("/nodes", s.NodeHandler.listNodesHandler).Methods("GET")
	s.router.HandleFunc("/nodes", s.NodeHandler.createNodeHandler).Methods("POST")
	s.router.HandleFunc("/nodes/{id}", s.NodeHandler.deleteNodeHandler).Methods("DELETE")
	s.router.HandleFunc("/deployments", s.DeploymentHandler.listDeploymentsHandler).Methods("GET")
	s.router.HandleFunc("/deployments/{name}", s.DeploymentHandler.deleteDeploymentHandler).Methods("DELETE")
	s.router.HandleFunc("/deployments/{name}/rollout-restart", s.DeploymentHandler.rolloutRestartDeploymentHandler).Methods("POST")
	s.router.HandleFunc("/deployments/{name}/scale", s.DeploymentHandler.scaleDeploymentHandler).Methods("POST")
	s.router.HandleFunc("/services", s.ServiceHandler.listServicesHandler).Methods("GET")
	s.router.HandleFunc("/services/{name}", s.ServiceHandler.deleteServiceHandler).Methods("DELETE")
	s.router.HandleFunc("/manifests", s.ManifestHandler.handleMadenResources).Methods("POST")
}

func (s *Server) Start() {
	go s.ChangeListener.WatchDeployments()
	go s.ChangeListener.WatchServices()
	go s.ChangeListener.WatchPodStatusChanges()

	server := &http.Server{
		Addr:    ":8080",
		Handler: s.router,
		ReadTimeout: 5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout: 1 * time.Minute,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("Failed to start server: %v", err)
		return
    }
	fmt.Println("Server is running on http://localhost:8080")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Maden API Server")
}
