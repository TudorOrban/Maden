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
	PodRestarted
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
	case "Restarted":
		*p = PodRestarted
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
	PersistentVolumeResource
	PersistentVolumeClaimResource
	DNSResource
)

func (r ResourceType) String() string {
	return [...]string{"Pod", "Node", "Deployment", "Service", "PersistentVolumeResource", "PersistentVolumeClaimResource", "DNSResource"}[r]
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

type ContainerStatus int

const (
	Created ContainerStatus = iota
	Running
	Paused
	Restarting
	Removing
	Exited
	Dead
)

func (c *ContainerStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	_, err := GetStatusFromString(s)
	if err != nil {
		return err
	}
	return nil
}

func GetStatusFromString(s string) (*ContainerStatus, error) {
	var c ContainerStatus
	switch s {
	case "created":
		c = Created
	case "running":
		c = Running
	case "paused":
		c = Paused
	case "restarting":
		c = Restarting
	case "removing":
		c = Removing
	case "exited":
		c = Exited
	case "dead":
		c = Dead
	default:
		return nil, fmt.Errorf("unknown container state: %s", s)
	}
	return &c, nil
}

func (c ContainerStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c ContainerStatus) String() string {
	return [...]string{"created", "running", "paused", "restarting", "removing", "exited", "dead"}[c]
}
