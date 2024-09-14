package main

import (
	"maden/pkg/apiserver"
	"maden/pkg/shared"

	"sync"
)

func main() {
	container := buildContainer()

	var wg sync.WaitGroup

	// DNS Server currently causes concurrency issues, so it is disabled for now
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()

	// 	err := container.Invoke(func(dnsServer *apiserver.DNSServer) {
	// 		dnsServer.StartDNSServer()
	// 	})
	// 	if err != nil {
	// 		shared.Log.Errorf("Failed to invoke DI container: %v", err)
	// 	}
	// }()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := container.Invoke(func(server *apiserver.Server) {
			server.Start()
		})
		if err != nil {
			shared.Log.Errorf("Failed to invoke DI container: %v", err)
		}
	}()

	wg.Wait()
}
