package main

import (
	"log"
	"maden/pkg/apiserver"
	"maden/pkg/etcd"

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
	container.Provide(etcd.NewEtcdDeploymentRepository)
	container.Provide(apiserver.NewDeploymentHandler)
	container.Provide(apiserver.NewServer)

	return container
}