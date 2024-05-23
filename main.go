package main

import (
	"maden/pkg/apiserver"

	"log"
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
