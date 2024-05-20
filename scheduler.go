package main


func schedulePod(pod *Pod) {
	mu.Lock()
	defer mu.Unlock()

	for i, node := range nodeDB {
		if shouldSchedulePod(&node, pod) {
			pod.NodeID = node.ID
			pod.Status = PodScheduled
			nodeDB[i].Used.CPU += pod.Resources.CPU
			nodeDB[i].Used.Memory += pod.Resources.Memory
			return
		}
	}

	pod.Status = PodPending
}

func shouldSchedulePod(node *Node, pod *Pod) bool {
	return node.Status == NodeReady && hasSufficentResources(node, &pod.Resources) &&
		matchesAffinity(node, pod) && matchesAntiAffinity(node, pod) &&
		matchesTolerations(node, pod)
}

func hasSufficentResources(node *Node, req *Resources) bool {
	availableCPU := node.Capacity.CPU - node.Used.CPU
	availableMemory := node.Capacity.Memory - node.Used.Memory
	return availableCPU >= req.CPU && availableMemory >= req.Memory
}

func matchesAffinity(node *Node, pod *Pod) bool {
	for key, val := range pod.Affinity {
		if nodeVal, ok := node.Labels[key]; !ok || nodeVal != val {
			return false
		}
	}
	return true
}

func matchesAntiAffinity(node *Node, pod *Pod) bool {
	for key, val := range pod.AntiAffinity {
		if nodeVal, ok := node.Labels[key]; ok && nodeVal == val {
			return false
		}
	}
	return true
}

func matchesTolerations(node *Node, pod *Pod) bool {
	for key, val := range node.Taints {
		if toVal, ok := pod.Tolerations[key]; !ok || toVal != val {
			return false
		}
	}
	return true
}
