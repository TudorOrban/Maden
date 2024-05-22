package controller

import "maden/pkg/shared"

type DeploymentController interface {
	HandleIncomingDeployment(deploymentSpec shared.DeploymentSpec) error
}