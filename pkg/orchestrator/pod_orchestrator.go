package orchestrator

import (
	"maden/pkg/etcd"
	"maden/pkg/madelet"
	"maden/pkg/scheduler"
	"maden/pkg/shared"

)


func OrchestratePodCreation(pod *shared.Pod) error {
	err := scheduler.SchedulePod(pod)
	if err != nil {
		return err
	}

	if err := etcd.CreatePod(pod); err != nil {
		return err
	}

	dockerRuntime, err := madelet.NewDockerRuntime()
	if err != nil {
		return err
	}

	podLifecycleManager := madelet.PodLifecycleManager{
		Runtime: dockerRuntime,
	}
	go podLifecycleManager.RunPod(pod)

	return nil
}