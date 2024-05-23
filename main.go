package main

import (
	"maden/pkg/apiserver"
	"maden/pkg/controller"
	"maden/pkg/etcd"
	"maden/pkg/orchestrator"

	"log"

	"go.uber.org/dig"
)

func main() {
	container := buildContainer()

	err := container.Invoke(func(server *apiserver.Server) {
		server.Start();
	})

	if err != nil {
		log.Fatalf("Failed to invoke container: %v", err)
	}

	etcd.InitEtcd()
	apiserver.InitAPIServer()
	
}

func buildContainer() *dig.Container {
	container := dig.New()

	container.Provide(etcd.NewEtcdClient)
	container.Provide(etcd.NewEtcdPodRepository)
	container.Provide(etcd.NewEtcdDeploymentRepository)
	container.Provide(controller.NewDefaultDeploymentController)
	container.Provide(controller.NewDefaultDeploymentUpdaterController)
	container.Provide(controller.NewEtcdChangeListener)
	container.Provide(orchestrator.NewDefaultPodOrchestrator)
	container.Provide(apiserver.NewPodHandler)
	container.Provide(apiserver.NewDeploymentHandler)
	container.Provide(apiserver.NewManifestHandler)
	container.Provide(apiserver.NewServer)

	return container
}