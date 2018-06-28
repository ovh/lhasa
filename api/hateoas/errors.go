package hateoas

import (
	"fmt"
	"sort"
)

// UnsupportedEntityError is raised when a repository method is called with a wrong type
type UnsupportedEntityError struct {
	Expected string
	Actual   string
}

// UnsupportedIndexError is raised when a impossible indexation was requested
type UnsupportedIndexError struct {
	Field     string
	Supported []string
}

// EntityDoesNotExistError is raised when a repository try to read a non existing entity
type EntityDoesNotExistError struct {
	EntityName string
	Criteria   map[string]interface{}
}

// InternalError all internal error
type InternalError struct {
	Message string
	Detail  string
}

type errorCreated string
type errorGone string

// Error implements error interface for UnsupportedEntityError
func (err UnsupportedEntityError) Error() string {
	return fmt.Sprintf("unsupported entity (expected: %s, actual: %s)", err.Expected, err.Actual)
}

// Error implements error interface for EntityDoesNotExistError
func (err EntityDoesNotExistError) Error() string {
	// sorting the criteria so the error messages are predictable (useful for testing)
	sortedCriteria := make([]string, 0)
	keys := make([]string, 0)
	for k := range err.Criteria {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sortedCriteria = append(sortedCriteria, fmt.Sprintf("%s=%s", k, err.Criteria[k]))
	}
	return fmt.Sprintf("entity %s%s does not exist", err.EntityName, sortedCriteria)
}

// Error implements error interface for UnsupportedIndexError
func (err UnsupportedIndexError) Error() string {
	return fmt.Sprintf("index by %s is not supported (one of %v)", err.Field, err.Supported)
}

// Error implements error interface
func (err *InternalError) Error() string {
	return err.Message + ":" + err.Detail
}

// Error implements error interface
func (err errorCreated) Error() string {
	return string(err)
}

// Error implements error interface
func (err errorGone) Error() string {
	return string(err)
}

// ErrorCreated is raised when no error occurs but a resource has been created (tonic single-status code workaround)
var ErrorCreated = errorCreated("created")

// ErrorGone is raised when a former resource has been requested but no longer exist
var ErrorGone = errorGone("gone")

//NewUnsupportedEntityError is a helper to create an UnsupportedEntityError
func NewUnsupportedEntityError(expected, actual interface{}) error {
	return UnsupportedEntityError{
		Expected: fmt.Sprintf("%T", expected),
		Actual:   fmt.Sprintf("%T", actual),
	}
}

// NewEntityDoesNotExistError is a helper to create an EntityDoesNotExistError
func NewEntityDoesNotExistError(entity interface{}, criteria map[string]interface{}) error {
	return EntityDoesNotExistError{
		EntityName: fmt.Sprintf("%T", entity),
		Criteria:   criteria,
	}
}

// NewUnsupportedIndexError is a helper to create an UnsupportedIndexError
func NewUnsupportedIndexError(field string, supported ...string) error {
	return UnsupportedIndexError{
		Field:     field,
		Supported: supported,
	}
}

// IsEntityDoesNotExistError returns true if err is an EntityDoesNotExistError
func IsEntityDoesNotExistError(err error) bool {
	_, ok := err.(EntityDoesNotExistError)
	return ok
}
