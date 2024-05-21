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

// Deployments & Services
type MadenResource struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind string `json:"kind" yaml:"kind"`
	Spec interface{} `json:"spec" yaml:"spec"`
}

// - Deployments
type DeploymentSpec struct {
	Name string `json:"name" yaml:"name"`
    Replicas int `json:"replicas" yaml:"replicas"`
    Selector LabelSelector `json:"selector" yaml:"selector"`
    Template PodTemplate `json:"template" yaml:"template"`
}

type Deployment struct {
	ID string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
	Replicas int `json:"replicas" yaml:"replicas"`
    Selector LabelSelector `json:"selector" yaml:"selector"`
	Template PodTemplate `json:"template" yaml:"template"`
}

type LabelSelector struct {
	MatchLabels map[string]string `json:"matchLabels" yaml:"matchLabels"`
}

type PodTemplate struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`
	Spec PodSpec `json:"spec" yaml:"spec"`
}

type PodSpec struct {
	Containers []Container `json:"containers" yaml:"containers"`
}

type Metadata struct {
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type Container struct {
	Image string `json:"image" yaml:"image"`
	Ports []Port `json:"ports" yaml:"ports"`
}

type Port struct {
	ContainerPort int `json:"containerPort" yaml:"containerPort"`
}

// - Services
type ServiceSpec struct {
	Name string `json:"name" yaml:"name"`
	Selector map[string]string `json:"selector" yaml:"selector"`
	Ports []ServicePort `json:"ports" yaml:"ports"`
}

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