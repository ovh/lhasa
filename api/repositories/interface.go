package repositories

import (
	"fmt"
	"math"
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

// Repository defines a normal repository
type Repository interface {
	FindAll() (interface{}, error)
	FindByID(interface{}) (interface{}, error)
	FindBy(map[string]interface{}) (interface{}, error)
	FindOneBy(map[string]interface{}) (interface{}, error)
	Save(interface{}) error
	Remove(interface{}) error
}

// PageableRepository defines a repository that handles pagination
type PageableRepository interface {
	Repository
	FindAllPage(Pageable) (Page, error)
	FindPageBy(Pageable, map[string]interface{}) (Page, error)
}

// TruncatableRepository defines a repository that handles truncation
type TruncatableRepository interface {
	Repository
	Truncate() error
}

// Pageable defines pagination criterias
type Pageable struct {
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
	Sort      string `json:"sort" form:"sort"`
	IndexedBy string `json:"indexedBy" form:"indexedBy"`
}

// Page defines a page
type Page struct {
	Content       interface{}
	TotalElements int      `json:"totalElements"`
	IDs           []string `json:"ids,omitempty"`
	Pageable      Pageable
}

// PagedResources defines a REST representation of paged contents
type PagedResources struct {
	Content      interface{}  `json:"content"`
	PageMetadata pageMetadata `json:"pageMetadata"`
	Links        []pageLink   `json:"_links"`
}

type pageMetadata struct {
	TotalElements int      `json:"totalElements"`
	TotalPages    int      `json:"totalPages"`
	Size          int      `json:"size"`
	Number        int      `json:"number"`
	IndexedBy     string   `json:"indexedBy,omitempty"`
	IDs           []string `json:"ids,omitempty"`
}

type pageLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// ToResources convert receiver page to a hatoas representation
func (page Page) ToResources(baseURL string) PagedResources {
	sortParam := ""
	if page.Pageable.Sort != "" {
		sortParam = "&sort=" + page.Pageable.Sort
	}
	links := []pageLink{
		{Href: fmt.Sprintf("%s?page=%d&size=%d%s", baseURL, page.Pageable.Page, page.Pageable.Size, sortParam), Rel: "self"},
	}
	totalPages := int(math.Ceil(float64(page.TotalElements) / float64(page.Pageable.Size)))
	if page.Pageable.Page < totalPages-1 {
		links = append(links, pageLink{Href: fmt.Sprintf("%s?page=%d&size=%d%s", baseURL, page.Pageable.Page+1, page.Pageable.Size, sortParam), Rel: "next"})
	}
	if page.Pageable.Page > 0 {
		links = append(links, pageLink{Href: fmt.Sprintf("%s?page=%d&size=%d%s", baseURL, page.Pageable.Page-1, page.Pageable.Size, sortParam), Rel: "prev"})
	}
	metadata := pageMetadata{
		Size:          page.Pageable.Size,
		Number:        page.Pageable.Page,
		TotalElements: page.TotalElements,
		TotalPages:    totalPages,
		IndexedBy:     page.Pageable.IndexedBy,
	}
	if len(page.IDs) > 0 {
		metadata.IDs = page.IDs
	}
	return PagedResources{
		Content:      page.Content,
		PageMetadata: metadata,
		Links:        links,
	}
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
