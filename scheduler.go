package main

func schedulePod(pod *Pod) {
	mu.Lock()
	defer mu.Unlock()

	// Round-robin scheduling
	if len(nodeDB) > 0 {
		for _, node := range nodeDB {
			if node.Status == NodeReady {
				pod.NodeID = node.ID
				pod.Status = PodScheduled
				return
			}
		}
	}

	pod.Status = PodPending
}