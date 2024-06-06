package controller

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"
	"maden/pkg/orchestrator"

	"fmt"
)

type DefaultServiceController struct {
	Repo etcd.ServiceRepository
	SvcOrchestrator orchestrator.ServiceOrchestrator
}

func NewDefaultServiceController(
	repo etcd.ServiceRepository,
	svcOrchestrator orchestrator.ServiceOrchestrator,
) ServiceController {
	return &DefaultServiceController{Repo: repo, SvcOrchestrator: svcOrchestrator}
}

func (c *DefaultServiceController) HandleIncomingService(serviceSpec shared.ServiceSpec) error {
	existingService, err := c.Repo.GetServiceByName(serviceSpec.Name)
	if err != nil {
		return c.SvcOrchestrator.OrchestrateServiceCreation(serviceSpec)
	}

	if existingService != nil && needsServiceUpdate(serviceSpec, existingService) {
		c.SvcOrchestrator.OrchestrateServiceUpdate(*existingService, serviceSpec)
	}

	fmt.Println("No update required for service: ", serviceSpec.Name)
	return nil
}


func needsServiceUpdate(spec shared.ServiceSpec, existing *shared.Service) bool {
	return !areMapsEqual(spec.Selector, existing.Selector) || 
	!arePortsEqual(spec.Ports, existing.Ports)
}


// Comparisons
func areMapsEqual(a, b map[string]string) bool {
    if len(a) != len(b) {
        return false
    }
    for key, valA := range a {
        if valB, ok := b[key]; !ok || valA != valB {
            return false
        }
    }
    return true
}

func arePortsEqual(a, b []shared.ServicePort) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i].Port != b[i].Port || a[i].TargetPort != b[i].TargetPort {
            return false
        }
    }
    return true
}
