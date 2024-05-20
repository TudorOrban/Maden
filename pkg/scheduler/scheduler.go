package scheduler

import (
	"maden/pkg/etcd"
	"maden/pkg/shared"
)

func SchedulePod(pod *shared.Pod) error {
	nodes, err := etcd.ListNodes()
	if err != nil {
		return err
	}

	etcd.Mu.Lock()
	defer etcd.Mu.Unlock()

	scheduled := false
	for i, node := range nodes {
		if shouldSchedulePod(&node, pod) {
			pod.NodeID = node.ID
			pod.Status = shared.PodScheduled
			nodes[i].Used.CPU += pod.Resources.CPU
			nodes[i].Used.Memory += pod.Resources.Memory

			if err := etcd.UpdateNode(&nodes[i]); err != nil {
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
	return node.Status == shared.NodeReady && hasSufficentResources(node, &pod.Resources) &&
		matchesAffinity(node, pod) && matchesAntiAffinity(node, pod) &&
		matchesTolerations(node, pod)
}

func hasSufficentResources(node *shared.Node, req *shared.Resources) bool {
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