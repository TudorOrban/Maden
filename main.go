package main

import (
	"maden/pkg/apiserver"
	"maden/pkg/controller"
	"maden/pkg/etcd"
	"maden/pkg/madelet"
	"maden/pkg/orchestrator"
	"maden/pkg/scheduler"

	"log"

	"go.uber.org/dig"
)

func main() {
	container := buildContainer()

	err := container.Invoke(func(server *apiserver.Server) {
		server.Start();
	})

	if err != nil {
		log.Fatalf("Failed to invoke DI container: %v", err)
	}	
}

func buildContainer() *dig.Container {
	container := dig.New()

	container.Provide(etcd.NewEtcdClient)
	container.Provide(madelet.NewDockerClient)
	container.Provide(etcd.NewEtcdPodRepository)
	container.Provide(etcd.NewEtcdNodeRepository)
	container.Provide(etcd.NewEtcdDeploymentRepository)
	container.Provide(etcd.NewEtcdServiceRepository)
	container.Provide(madelet.NewContainerRuntimeInterface)
	container.Provide(scheduler.NewPodScheduler)
	container.Provide(controller.NewDefaultDeploymentController)
	container.Provide(controller.NewDefaultDeploymentUpdaterController)
	container.Provide(controller.NewDefaultServiceController)
	container.Provide(controller.NewDefaultServiceUpdaterController)
	container.Provide(controller.NewEtcdChangeListener)
	container.Provide(madelet.NewPodLifecycleManager)
	container.Provide(orchestrator.NewDefaultPodOrchestrator)
	container.Provide(apiserver.NewPodHandler)
	container.Provide(apiserver.NewNodeHandler)
	container.Provide(apiserver.NewDeploymentHandler)
	container.Provide(apiserver.NewServiceHandler)
	container.Provide(apiserver.NewManifestHandler)
	container.Provide(apiserver.NewServer)

	return container
}