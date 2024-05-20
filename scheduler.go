package main

func schedulePod(pod *Pod) {
	mu.Lock()
	defer mu.Unlock()

	for i, node := range nodeDB {
		if node.Status == NodeReady && hasSufficentResources(&node, &pod.Resources) &&
		   matchesAffinity(&node, pod) && matchesAntiAffinity(&node, pod) {
			pod.NodeID = node.ID
			pod.Status = PodScheduled
			nodeDB[i].Used.CPU += pod.Resources.CPU
			nodeDB[i].Used.Memory += pod.Resources.Memory
			return
		}
	}

	pod.Status = PodPending
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