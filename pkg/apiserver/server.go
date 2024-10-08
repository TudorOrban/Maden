package apiserver

import (
	"maden/pkg/controller"
	"maden/pkg/shared"

	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	router            *mux.Router
	PodHandler        *PodHandler
	NodeHandler       *NodeHandler
	DeploymentHandler *DeploymentHandler
	ServiceHandler    *ServiceHandler
	PersistentVolumeHandler *PersistentVolumeHandler
	PermanentVolumeClaimHandler *PersistentVolumeClaimHandler
	ManifestHandler   *ManifestHandler

	ChangeListener *controller.EtcdChangeListener
}

func NewServer(
	podHandler *PodHandler,
	nodeHandler *NodeHandler,
	deploymentHandler *DeploymentHandler,
	serviceHandler *ServiceHandler,
	persistentVolumeHandler *PersistentVolumeHandler,
	persistentVolumeClaimHandler *PersistentVolumeClaimHandler,
	manifestHandler *ManifestHandler,
	changeListener *controller.EtcdChangeListener,
) *Server {
	s := &Server{
		router:            mux.NewRouter(),
		PodHandler:        podHandler,
		NodeHandler:       nodeHandler,
		DeploymentHandler: deploymentHandler,
		ServiceHandler:    serviceHandler,
		PersistentVolumeHandler: persistentVolumeHandler,
		PermanentVolumeClaimHandler: persistentVolumeClaimHandler,
		ManifestHandler:   manifestHandler,
		ChangeListener:    changeListener,
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
	s.router.HandleFunc("/persistent-volumes", s.PersistentVolumeHandler.listPersistentVolumesHandler).Methods("GET")
	s.router.HandleFunc("/persistent-volumes/{id}", s.PersistentVolumeHandler.deletePersistentVolumeHandler).Methods("DELETE")
	s.router.HandleFunc("/persistent-volume-claims", s.PermanentVolumeClaimHandler.listPersistentVolumeClaimsHandler).Methods("GET")
	s.router.HandleFunc("/persistent-volume-claims/{id}", s.PermanentVolumeClaimHandler.deletePersistentVolumeClaimHandler).Methods("DELETE")
	s.router.HandleFunc("/manifests", s.ManifestHandler.handleMadenResources).Methods("POST")
}

func (s *Server) Start() {
	go s.ChangeListener.WatchDeployments()
	go s.ChangeListener.WatchServices()
	go s.ChangeListener.WatchPodStatusChanges()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      s.router,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  1 * time.Minute,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		shared.Log.Errorf("Failed to start server: %v", err)
		return
	}
	fmt.Println("Server is running on http://localhost:8080")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Maden API Server")
}
