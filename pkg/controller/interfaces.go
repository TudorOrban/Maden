package controller

import (
	"maden/pkg/shared"

	"go.etcd.io/etcd/api/v3/mvccpb"
)

type DeploymentController interface {
	HandleIncomingDeployment(deploymentSpec shared.DeploymentSpec) error
}

type DeploymentUpdaterController interface {
	HandleDeploymentCreate(kv *mvccpb.KeyValue)
	HandleDeploymentUpdate(oldKv *mvccpb.KeyValue, newKv *mvccpb.KeyValue)
	HandleDeploymentDelete(kv *mvccpb.KeyValue)
}