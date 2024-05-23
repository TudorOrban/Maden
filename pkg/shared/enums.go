package shared

import (
	"encoding/json"
	"fmt"
)

type NodeStatus int

const (
	NodeReady NodeStatus = iota
	NodeNotReady
	NodeOffline
)

func (n *NodeStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "Ready":
		*n = NodeReady
	case "NotReady":
		*n = NodeNotReady
	case "Offline":
		*n = NodeOffline
	default:
		return fmt.Errorf("unknown node status: %s", s)
	}
	return nil
}

func (n NodeStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

func (n NodeStatus) String() string {
	return [...]string{"Ready", "NotReady", "Offline"}[n]
}

type PodStatus int

const (
	PodPending PodStatus = iota
	PodScheduled
	PodContainerCreating
	PodRunning
	PodFailed
)

func (p *PodStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "Pending":
		*p = PodPending
	case "Scheduled":
		*p = PodScheduled
	case "ContainerCreating":
		*p = PodContainerCreating
	case "Running":
		*p = PodRunning
	case "Failed":
		*p = PodFailed
	default:
		return fmt.Errorf("unknown pod status: %s", s)
	}
	return nil
}

func (p PodStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p PodStatus) String() string {
	return [...]string{"Pending", "Scheduled", "ContainerCreating", "Running", "Failed"}[p]
}


type ResourceType int

const (
	PodResource ResourceType = iota
	NodeResource
	DeploymentResource
	ServiceResource
)

func (r ResourceType) String() string {
	return [...]string{"Pod", "Node", "Deployment", "Service"}[r]
}

type RestartPolicy int

const (
	RestartAlways RestartPolicy = iota
	RestartOnFailure
	RestartNever
)

func (r *RestartPolicy) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "Always":
		*r = RestartAlways
	case "OnFailure":
		*r = RestartOnFailure
	case "Never":
		*r = RestartNever
	default:
		return fmt.Errorf("unknown restart policy: %s", s)
	}
	return nil
}

func (r RestartPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r RestartPolicy) String() string {
	return [...]string{"Always", "OnFailure", "Never"}[r]
}