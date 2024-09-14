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
	HandleDeploymentRolloutRestart(deployment *shared.Deployment) error
}

type ServiceController interface {
	HandleIncomingService(serviceSpec shared.ServiceSpec) error
}

type ServiceUpdaterController interface {
	HandleServiceCreate(kv *mvccpb.KeyValue)
	HandleServiceUpdate(prevKv *mvccpb.KeyValue, newKv *mvccpb.KeyValue)
	HandleServiceDelete(prevKv *mvccpb.KeyValue)
}

type PodUpdaterController interface {
	HandlePodUpdate(oldKv *mvccpb.KeyValue, newKv *mvccpb.KeyValue)
}

type PersistentVolumeController interface {
	HandleIncomingPersistentVolume(volumeSpec shared.PersistentVolumeSpec) error
}

type PersistentVolumeClaimController interface {
	HandleIncomingPersistentVolumeClaim(volumeClaimSpec shared.PersistentVolumeClaimSpec) error
}