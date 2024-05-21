package shared

import (
	"fmt"
)

type ErrNotFound struct {
	ID string
	Name string
	ResourceType ResourceType
}

func (e *ErrNotFound) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("a %s with ID %s could not be found", e.ResourceType.String(), e.ID)
	} else if e.Name != "" {
		return fmt.Sprintf("a %s with name %s could not be found", e.ResourceType.String(), e.Name)
	}
	return fmt.Sprintf("resource %s could not be found", e.ResourceType.String())
}

type ErrDuplicateResource struct {
	ID string
	ResourceType ResourceType
}

func (e *ErrDuplicateResource) Error() string {
	return fmt.Sprintf("a %s with ID %s already exists", e.ResourceType.String(), e.ID)
}