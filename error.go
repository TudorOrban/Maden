package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("resource not found")

type ErrDuplicateResource struct {
	ID string
	ResourceType ResourceType
}

func (e *ErrDuplicateResource) Error() string {
	return fmt.Sprintf("a pod with ID %s already exists", e.ID)
}

type ResourceType int

const (
	PodResource ResourceType = iota
	NodeResource
)