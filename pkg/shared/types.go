package shared

// Nodes
type Node struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Status NodeStatus `json:"status"`
	Capacity Resources `json:"capacity"`
	Used Resources `json:"used"`
	Labels map[string]string `json:"labels"`
	Taints map[string]string `json:"taints"`
}

type NodeCapacity struct {
	CPU int `json:"cpu"`
	Memory int `json:"memory"`
}

// Pods
type Pod struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Status PodStatus `json:"status"`
	NodeID string `json:"nodeId"`
	Resources Resources `json:"resources"`
	Affinity map[string]string `json:"affinity"`
	AntiAffinity map[string]string `json:"antiAffinity"`
	Tolerations map[string]string `json:"tolerations"`
}

type Resources struct {
	CPU int `json:"cpu"`
	Memory int `json:"memory"` // in MB
}