package orchestrator

import (
	"maden/pkg/etcd"
	"maden/pkg/networking"
	"maden/pkg/shared"
)

type DefaultServiceOrchestrator struct {
	Repo      etcd.ServiceRepository
	DNSRepo   etcd.DNSRepository
	IPManager networking.IPManager
}

func NewDefaultServiceOrchestrator(
	repo etcd.ServiceRepository,
	dnsRepo etcd.DNSRepository,
	ipManager networking.IPManager,
) ServiceOrchestrator {
	return &DefaultServiceOrchestrator{Repo: repo, DNSRepo: dnsRepo, IPManager: ipManager}
}

func (o *DefaultServiceOrchestrator) OrchestrateServiceCreation(serviceSpec shared.ServiceSpec) error {
	service := transformToService(serviceSpec)

	ip, err := o.IPManager.AssignIP()
	if err != nil {
		return err
	}
	service.IP = ip

	if err := o.Repo.CreateService(&service); err != nil {
		return err
	}

	if err := o.DNSRepo.RegisterService(service.Name, service.IP); err != nil {
		return err
	}

	return nil
}

func transformToService(spec shared.ServiceSpec) shared.Service {
	id := shared.GenerateRandomString(10)
	service := shared.Service{
		ID:       id,
		Name:     spec.Name,
		Selector: spec.Selector,
		Ports:    spec.Ports,
	}
	return service
}

func (o *DefaultServiceOrchestrator) OrchestrateServiceUpdate(existingService shared.Service, serviceSpec shared.ServiceSpec) error {
	shared.Log.Infof("Updating service...")
	updatedService := updateExistingService(serviceSpec, &existingService)
	err := o.Repo.UpdateService(&updatedService)
	if err != nil {
		return err
	}

	return o.DNSRepo.RegisterService(updatedService.Name, updatedService.IP)
}

func updateExistingService(spec shared.ServiceSpec, existing *shared.Service) shared.Service {
	(*existing).Selector = spec.Selector
	(*existing).Ports = spec.Ports
	return *existing
}

func (o *DefaultServiceOrchestrator) OrchestrateServiceDeletion(serviceName string) error {
	service, err := o.Repo.GetServiceByName(serviceName)
	if err != nil {
		return err
	}

	if err := o.DNSRepo.DeregisterService(service.Name); err != nil {
		shared.Log.Errorf("failed to deregister service %s: %v", service.Name, err)
	}

	if err := o.IPManager.ReleaseIP(service.IP); err != nil {
		shared.Log.Errorf("failed to release IP %s: %v", service.IP, err)
	}

	return o.Repo.DeleteService(service.Name)
}
