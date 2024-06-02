package main

import (
	"maden/pkg/apiserver"
	"maden/pkg/networking"
	"sync"

	"log"
)


func main() {
	container := buildContainer()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := container.Invoke(func(dnsServer *networking.DNSServer) {
			dnsServer.StartDNSServer();
		})
		if err != nil {
			log.Fatalf("Failed to invoke DI container: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := container.Invoke(func(server *apiserver.Server) {
			server.Start();
		})
		if err != nil {
			log.Fatalf("Failed to invoke DI container: %v", err)
		}
	}()
	
	wg.Wait()
}
