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

// Deployments
type DeploymentSpec struct {
    APIVersion string `yaml:"version"`
    Deployment Deployment `yaml:"deployment"`
}

type Deployment struct {
	ID string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	Replicas int `json:"replicas" yaml:"replicas"`
	Selector map[string]string `json:"selector" yaml:"selector"`
	Template PodTemplate `json:"template" yaml:"template"`
}

type PodTemplate struct {
	Containers []Container `json:"containers" yaml:"containers"`
}

type Container struct {
	Image string `json:"image" yaml:"image"`
	Ports []Port `json:"ports" yaml:"ports"`
}

type Port struct {
	ContainerPort int `json:"containerPort" yaml:"containerPort"`
}

// Services
type Service struct {
	ID string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	Selector map[string]string `json:"selector" yaml:"selector"`
	Ports []ServicePort `json:"ports" yaml:"ports"`
}

type ServicePort struct {
	Port int `json:"port" yaml:"port"`
	TargetPort int `json:"targetPort" yaml:"targetPort"`
}