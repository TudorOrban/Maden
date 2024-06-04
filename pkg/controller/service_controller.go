package controller

import (
	"maden/pkg/etcd"
	"maden/pkg/networking"
	"maden/pkg/shared"

	"fmt"
)

type DefaultServiceController struct {
	Repo etcd.ServiceRepository
	DNSRepo etcd.DNSRepository
	IPManager networking.IPManager
}

func NewDefaultServiceController(
	repo etcd.ServiceRepository,
	dnsRepo etcd.DNSRepository,
	ipManager networking.IPManager,
) ServiceController {
	return &DefaultServiceController{Repo: repo, DNSRepo: dnsRepo, IPManager: ipManager}
}

func (c *DefaultServiceController) HandleIncomingService(serviceSpec shared.ServiceSpec) error {
	existingService, err := c.Repo.GetServiceByName(serviceSpec.Name)
	if err != nil {
		if _, ok := err.(*shared.ErrNotFound); ok {
			fmt.Println("Creating service")
			service := transformToService(serviceSpec)
			
			ip, err := c.IPManager.AssignIP()
			if err != nil {
				return err
			}
			service.IP = ip

			err = c.Repo.CreateService(&service)
			if err != nil {
				return err
			}

			c.DNSRepo.RegisterService(service.Name, service.IP)
		} else {
			return err
		}
	}

	if needsServiceUpdate(serviceSpec, existingService) {
		fmt.Println("Updating service")
		updatedService := updateExistingService(serviceSpec, existingService)
		err = c.Repo.UpdateService(&updatedService)
		if err != nil {
			return err
		}
		
		c.DNSRepo.RegisterService(updatedService.Name, updatedService.IP)
	}

	fmt.Println("No update required for service: ", serviceSpec.Name)
	return nil
}

func transformToService(spec shared.ServiceSpec) shared.Service {
	id := shared.GenerateRandomString(10)
	service := shared.Service{
		ID: id,
		Name: spec.Name,
		Selector: spec.Selector,
		Ports: spec.Ports,
	}
	return service
}

func needsServiceUpdate(spec shared.ServiceSpec, existing *shared.Service) bool {
	return !areMapsEqual(spec.Selector, existing.Selector) || 
	!arePortsEqual(spec.Ports, existing.Ports)
}

func updateExistingService(spec shared.ServiceSpec, existing *shared.Service) shared.Service {
	(*existing).Selector = spec.Selector
	(*existing).Ports = spec.Ports
	return *existing
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
