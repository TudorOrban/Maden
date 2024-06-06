package orchestrator

import (
	"maden/pkg/etcd"
	"maden/pkg/networking"

	"log"
)

type DefaultServiceOrchestrator struct {
	Repo etcd.ServiceRepository
	DNSRepo etcd.DNSRepository
	IPManager networking.IPManager
}

func NewDefaultServiceOrchestrator(
	repo etcd.ServiceRepository,
	dnsRepo etcd.DNSRepository,
	ipManager networking.IPManager,
) ServiceOrchestrator {
	return &DefaultServiceOrchestrator{Repo: repo, DNSRepo: dnsRepo, IPManager: ipManager}
}

func (o *DefaultServiceOrchestrator) OrchestrateServiceDeletion(serviceName string) error {
	service, err := o.Repo.GetServiceByName(serviceName)
	if err != nil {
		return err
	}

	if err := o.DNSRepo.DeregisterService(service.Name); err != nil {
		log.Printf("failed to deregister service %s: %v", service.Name, err)
	}

	if err := o.IPManager.ReleaseIP(service.IP); err != nil {
		log.Printf("failed to release IP %s: %v", service.IP, err)
	}

	return o.Repo.DeleteService(service.Name)
}