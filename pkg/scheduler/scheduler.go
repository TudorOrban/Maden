package scheduler

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"
)

type PodScheduler struct {
	Repo etcd.NodeRepository
}

func NewPodScheduler(repo etcd.NodeRepository) Scheduler {
	return &PodScheduler{Repo: repo}
}

func (s *PodScheduler) SchedulePod(pod *shared.Pod) error {
	nodes, err := s.Repo.ListNodes()
	if err != nil {
		return err
	}

	scheduled := false
	for i, node := range nodes {
		if shouldSchedulePod(&node, pod) {
			pod.NodeID = node.ID
			pod.Status = shared.PodScheduled
			nodes[i].Used.CPU += pod.Resources.CPU
			nodes[i].Used.Memory += pod.Resources.Memory

			if err := s.Repo.UpdateNode(&nodes[i]); err != nil {
				return err
			}

			scheduled = true
			break
		}
	}

	if !scheduled {
		pod.Status = shared.PodPending
	}

	return nil
}

func shouldSchedulePod(node *shared.Node, pod *shared.Pod) bool {
	return node.Status == shared.NodeReady && hasSufficientResources(node, &pod.Resources) &&
		matchesAffinity(node, pod) && matchesAntiAffinity(node, pod) &&
		matchesTolerations(node, pod)
}

func hasSufficientResources(node *shared.Node, req *shared.Resources) bool {
	availableCPU := node.Capacity.CPU - node.Used.CPU
	availableMemory := node.Capacity.Memory - node.Used.Memory
	return availableCPU >= req.CPU && availableMemory >= req.Memory
}

func matchesAffinity(node *shared.Node, pod *shared.Pod) bool {
	for key, val := range pod.Affinity {
		if nodeVal, ok := node.Labels[key]; !ok || nodeVal != val {
			return false
		}
	}
	return true
}

func matchesAntiAffinity(node *shared.Node, pod *shared.Pod) bool {
	for key, val := range pod.AntiAffinity {
		if nodeVal, ok := node.Labels[key]; ok && nodeVal == val {
			return false
		}
	}
	return true
}

func matchesTolerations(node *shared.Node, pod *shared.Pod) bool {
	for key, val := range node.Taints {
		if toVal, ok := pod.Tolerations[key]; !ok || toVal != val {
			return false
		}
	}
	return true
}