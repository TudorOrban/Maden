package etcd

import "maden/pkg/shared"

type DeploymentRepository interface {
	ListDeployments() ([]shared.Deployment, error)
	GetDeploymentByName(deploymentName string) (*shared.Deployment, error)
	CreateDeployment(deployment *shared.Deployment) error
	UpdateDeployment(deployment *shared.Deployment) error
	DeleteDeployment(deploymentName string) error
}

type PodRepository interface {
	ListPods() ([]shared.Pod, error)
	GetPodsByDeploymentID(deploymentID string) ([]shared.Pod, error)
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