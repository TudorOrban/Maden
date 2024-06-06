package main

import (
	"maden/pkg/apiserver"
	"maden/pkg/controller"
	"maden/pkg/etcd"
	"maden/pkg/madelet"
	"maden/pkg/networking"
	"maden/pkg/orchestrator"
	"maden/pkg/scheduler"

	"go.uber.org/dig"
)

func buildContainer() *dig.Container {
	container := dig.New()

	container.Provide(etcd.NewEtcdClient)
	container.Provide(etcd.ProvideEtcdClient)
	container.Provide(madelet.NewClient)
	container.Provide(madelet.NewDockerClient)
	container.Provide(etcd.NewEtcdPodRepository)
	container.Provide(etcd.NewEtcdNodeRepository)
	container.Provide(etcd.NewEtcdDeploymentRepository)
	container.Provide(etcd.NewEtcdServiceRepository)
	container.Provide(etcd.NewEtcdTransactionRepository)
	container.Provide(etcd.NewEtcdDNSRepository)
	container.Provide(madelet.NewContainerRuntimeInterface)
	container.Provide(scheduler.NewPodScheduler)
	container.Provide(controller.NewDefaultDeploymentController)
	container.Provide(controller.NewDefaultDeploymentUpdaterController)
	container.Provide(controller.NewDefaultServiceController)
	container.Provide(controller.NewDefaultServiceUpdaterController)
	container.Provide(controller.NewDefaultPodUpdaterController)
	container.Provide(controller.NewEtcdChangeListener)
	container.Provide(madelet.NewPodLifecycleManager)
	container.Provide(orchestrator.NewDefaultPodOrchestrator)
	container.Provide(orchestrator.NewDefaultServiceOrchestrator)
	container.Provide(func() networking.IPManager {
		return networking.NewSimpleIPManager()
	})
	container.Provide(apiserver.NewPodHandler)
	container.Provide(apiserver.NewNodeHandler)
	container.Provide(apiserver.NewDeploymentHandler)
	container.Provide(apiserver.NewServiceHandler)
	container.Provide(apiserver.NewManifestHandler)
	container.Provide(apiserver.NewDNSHandler)
	container.Provide(apiserver.NewServer)
	container.Provide(apiserver.NewDNSServer)

	return container
}