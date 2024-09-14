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
	DeploymentID string `json:"deploymentId"`
	Status PodStatus `json:"status"`
	NodeID string `json:"nodeId"`
	Containers []Container `json:"containers"`
	Resources Resources `json:"resources"`
	Affinity map[string]string `json:"affinity"`
	AntiAffinity map[string]string `json:"antiAffinity"`
	Tolerations map[string]string `json:"tolerations"`
	RestartPolicy RestartPolicy `json:"restartPolicy" yaml:"restartPolicy"`
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
	Resources Resources `json:"resources" yaml:"resources"`
	Affinity map[string]string `json:"affinity" yaml:"affinity"`
	AntiAffinity map[string]string `json:"antiAffinity" yaml:"antiAffinity"`
	Tolerations map[string]string `json:"tolerations" yaml:"tolerations"`
	RestartPolicy RestartPolicy `json:"restartPolicy" yaml:"restartPolicy"`
}

type Metadata struct {
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type Container struct {
	ID string `json:"containerId"`
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
	IP string `json:"ip" yaml:"ip"`
}

type ServicePort struct {
	Port int `json:"port" yaml:"port"`
	TargetPort int `json:"targetPort" yaml:"targetPort"`
}

// Persistent Volumes
type PersistentVolumeSpec struct {
    Capacity map[string]string `json:"capacity" yaml:"capacity"`
    AccessModes []string `json:"accessModes" yaml:"accessModes"`
    PersistentVolumeReclaimPolicy string `json:"reclaimPolicy" yaml:"reclaimPolicy"`
    StorageClassName string `json:"storageClassName" yaml:"storageClassName"`
    MountOptions []string `json:"mountOptions" yaml:"mountOptions"`
}

type PersistentVolume struct {
	ID string `json:"id" yaml:"id"`
    Metadata Metadata `json:"metadata" yaml:"metadata"`
    Spec PersistentVolumeSpec `json:"spec" yaml:"spec"`
}

type PersistentVolumeClaimSpec struct {
    AccessModes []string `json:"accessModes" yaml:"accessModes"`
    Resources Resources `json:"resources" yaml:"resources"`
    VolumeName string `json:"volumeName" yaml:"volumeName"`
}

type PersistentVolumeClaim struct {
	ID string `json:"id" yaml:"id"`
    Metadata Metadata `json:"metadata" yaml:"metadata"`
    Spec PersistentVolumeClaimSpec `json:"spec" yaml:"spec"`
}

// Other 
type ScaleRequest struct {
	Replicas int `json:"replicas"`
}