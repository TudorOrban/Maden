package controller

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"

	"fmt"
)

type DefaultDeploymentController struct {
	Repo etcd.DeploymentRepository
}

func NewDefaultDeploymentController(repo etcd.DeploymentRepository) DeploymentController {
	return &DefaultDeploymentController{Repo: repo}
}


func (c *DefaultDeploymentController) HandleIncomingDeployment(deploymentSpec shared.DeploymentSpec) error {
	existingDeployment, err := c.Repo.GetDeploymentByName(deploymentSpec.Name)
	if err != nil {
		if _, ok := err.(*shared.ErrNotFound); ok {
			fmt.Println("Creating deployment")
			deployment := transformToDeployment(deploymentSpec)
			return c.Repo.CreateDeployment(&deployment)
		} else {
			return err
		}
	}

	if needsDeploymentUpdate(deploymentSpec, existingDeployment) {
		fmt.Println("Updating deployment")
		existingDeployment := updateExistingDeployment(deploymentSpec, existingDeployment)
		return c.Repo.UpdateDeployment(&existingDeployment)
	}

	fmt.Println("No update required for deployment: ", deploymentSpec.Name)
	return nil
}

func transformToDeployment(spec shared.DeploymentSpec) shared.Deployment {
	id := shared.GenerateRandomString(10)
	deployment := shared.Deployment{
		ID: id,
		Name: spec.Name,
		Replicas: spec.Replicas,
		Selector: spec.Selector,
		Template: spec.Template,
	}
	return deployment
}

func needsDeploymentUpdate(spec shared.DeploymentSpec, existing *shared.Deployment) bool {
	return spec.Replicas != existing.Replicas || 
	!areSelectorsEqual(spec.Selector, existing.Selector) || 
	!arePodTemplatesEqual(spec.Template, existing.Template)
}

func updateExistingDeployment(spec shared.DeploymentSpec, existing *shared.Deployment) shared.Deployment {
	(*existing).Replicas = spec.Replicas
	(*existing).Selector = spec.Selector
	(*existing).Template = spec.Template
	return *existing
}

// Comparisons
func areSelectorsEqual(a, b shared.LabelSelector) bool {
    return areMapsEqual(a.MatchLabels, b.MatchLabels)
}

func arePodTemplatesEqual(a, b shared.PodTemplate) bool {
    if !areMapsEqual(a.Metadata.Labels, b.Metadata.Labels) {
        return false
    }
    return arePodSpecsEqual(a.Spec, b.Spec)
}

func arePodSpecsEqual(a, b shared.PodSpec) bool {
    if len(a.Containers) != len(b.Containers) {
        return false
    }
    for i := range a.Containers {
        if !areContainersEqual(a.Containers[i], b.Containers[i]) {
            return false
        }
    }
    return true
}

func areContainersEqual(a, b shared.Container) bool {
    if a.Image != b.Image {
        return false
    }
    if len(a.Ports) != len(b.Ports) {
        return false
    }
    for i := range a.Ports {
        if a.Ports[i].ContainerPort != b.Ports[i].ContainerPort {
            return false
        }
    }
    return true
}