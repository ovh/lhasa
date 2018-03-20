package repositories

import "fmt"

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

// Error implements error interface for UnsupportedEntityError
func (err UnsupportedEntityError) Error() string {
	return fmt.Sprintf("unsupported entity (expected: %s, actual: %s)", err.Expected, err.Actual)
}

// Error implements error interface for EntityDoesNotExistError
func (err EntityDoesNotExistError) Error() string {
	return fmt.Sprintf("entity %s of type %s does not exist", err.Criteria, err.EntityName)
}

// Error implements error interface for UnsupportedIndexError
func (err UnsupportedIndexError) Error() string {
	return fmt.Sprintf("index by %s is not supported (one of %v)", err.Field, err.Supported)
}

//NewUnsupportedEntityError is a helper to create an UnsupportedEntityError
func NewUnsupportedEntityError(expected, actual interface{}) error {
	return UnsupportedEntityError{
		Expected: fmt.Sprintf("%T", expected),
		Actual:   fmt.Sprintf("%T", actual),
	}
}

// NewEntityDoesNotExistError is a helper to create an EntityDoesNotExistError
func NewEntityDoesNotExistError(entity interface{}, criterias map[string]interface{}) error {
	return EntityDoesNotExistError{
		EntityName: fmt.Sprintf("%T", entity),
		Criteria:   criterias,
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
