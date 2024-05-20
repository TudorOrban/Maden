package shared

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
	return fmt.Sprintf("a %s with ID %s already exists", e.ResourceType.String(), e.ID)
}

