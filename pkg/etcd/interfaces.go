package etcd

import (
	"context"
	"maden/pkg/shared"
)

type PodRepository interface {
	ListPods() ([]shared.Pod, error)
	GetPodsByDeploymentID(deploymentID string) ([]shared.Pod, error)
	GetPodByID(podID string) (*shared.Pod, error)
	CreatePod(pod *shared.Pod) error
	UpdatePod(pod *shared.Pod) error
	DeletePod(podID string) error
}

type NodeRepository interface {
	ListNodes() ([]shared.Node, error)
	CreateNode(node *shared.Node) error
	UpdateNode(node *shared.Node) error
	DeleteNode(nodeName string) error
}

type DeploymentRepository interface {
	ListDeployments() ([]shared.Deployment, error)
	GetDeploymentByName(deploymentName string) (*shared.Deployment, error)
	CreateDeployment(deployment *shared.Deployment) error
	UpdateDeployment(deployment *shared.Deployment) error
	DeleteDeployment(deploymentName string) error
}

type ServiceRepository interface {
	ListServices() ([]shared.Service, error)
	GetServiceByName(serviceName string) (*shared.Service, error)
	CreateService(service *shared.Service) error
	UpdateService(service *shared.Service) error
	DeleteService(serviceName string) error
}

type PersistentVolumeRepository interface {
	ListPersistentVolumes() ([]shared.PersistentVolume, error)
	GetPersistentVolumeByID(persistentVolumeID string) (*shared.PersistentVolume, error)
	CreatePersistentVolume(volume *shared.PersistentVolume) error
	UpdatePersistentVolume(volume *shared.PersistentVolume) error
	DeletePersistentVolume(volumeName string) error
}

type PersistentVolumeClaimRepository interface {
	ListPersistentVolumeClaims() ([]shared.PersistentVolumeClaim, error)
	GetPersistentVolumeClaimByID(persistentVolumeClaimID string) (*shared.PersistentVolumeClaim, error)
	CreatePersistentVolumeClaim(volumeClaim *shared.PersistentVolumeClaim) error
	UpdatePersistentVolumeClaim(volumeClaim *shared.PersistentVolumeClaim) error
	DeletePersistentVolumeClaim(volumeClaimName string) error
}

type Transactioner interface {
	PerformTransaction(ctx context.Context, key string, value string, resourceType shared.ResourceType) error
}

type DNSRepository interface {
	RegisterService(serviceName string, serviceIP string) error
	DeregisterService(serviceName string) error
	ResolveService(serviceName string) (string, error)
}
