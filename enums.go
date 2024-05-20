package main

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
	PodRunning
	PodScheduled
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
	case "Running":
		*p = PodRunning
	case "Scheduled":
		*p = PodScheduled
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
	return [...]string{"Pending", "Running", "Scheduled", "Failed"}[p]
}
