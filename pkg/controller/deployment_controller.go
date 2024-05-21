package controller

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"
	
	"fmt"
)

func HandleIncomingDeployment(deploymentSpec shared.DeploymentSpec) error {
	existingDeployment, err := etcd.GetDeploymentByName(deploymentSpec.Name)
	if err != nil {
		if _, ok := err.(*shared.ErrNotFound); ok {
			fmt.Println("Creating deployment")
			deployment := transformToDeployment(deploymentSpec)
			return etcd.CreateDeployment(&deployment)
		} else {
			return err
		}
	}

	if needsUpdate(deploymentSpec, existingDeployment) {
		fmt.Println("Updating deployment")
		existingDeployment := updateExistingDeployment(deploymentSpec, existingDeployment)
		return etcd.UpdateDeployment(&existingDeployment)
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

func needsUpdate(spec shared.DeploymentSpec, existing *shared.Deployment) bool {
	return spec.Replicas != existing.Replicas
}

func updateExistingDeployment(spec shared.DeploymentSpec, existing *shared.Deployment) shared.Deployment {
	(*existing).Replicas = spec.Replicas
	(*existing).Selector = spec.Selector
	(*existing).Template = spec.Template
	return *existing
}
