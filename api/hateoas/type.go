package hateoas

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

// Sorting directions
const (
	directionAsc  = "asc"
	directionDesc = "desc"
)

// Entity defines and identifiable REST entity
type Entity interface {
	GetID() string
	SetID(string) error
}

// SoftDeletableEntity defines a soft deletable REST entity
type SoftDeletableEntity interface {
	Entity
	GetDeletedAt() *time.Time
}

// BaseRepository defines a normal repository
type BaseRepository interface {
	GetType() reflect.Type
}

// ListableRepository defines a repository where one can list entities
type ListableRepository interface {
	BaseRepository
	FindBy(map[string]interface{}) (interface{}, error)
	FindOneBy(map[string]interface{}) (Entity, error)
}

// SavableRepository defines a repository where one can persist or remove entities
type SavableRepository interface {
	ListableRepository
	Save(Entity) error
	Remove(interface{}) error
}

// PageableRepository defines a repository that handles pagination
type PageableRepository interface {
	BaseRepository
	FindPageBy(Pageable, map[string]interface{}) (Page, error)
}

// TruncatableRepository defines a repository that handles truncation
type TruncatableRepository interface {
	BaseRepository
	Truncate() error
}

// SoftDeletableRepository defines a repository that handles soft delete
type SoftDeletableRepository interface {
	BaseRepository
	FindOneByUnscoped(criteria map[string]interface{}) (SoftDeletableEntity, error)
}

// Pageable defines pagination criteria
type Pageable struct {
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
	Sort      string `json:"sort" form:"sort"`
	IndexedBy string `json:"indexedBy" form:"indexedBy" default:""`
}

// Page defines a page
type Page struct {
	Content       interface{}
	TotalElements int      `json:"totalElements"`
	IDs           []string `json:"ids,omitempty"`
	BasePath      string   `json:"-"`
	Pageable      Pageable
}

// PagedResources defines a REST representation of paged contents
type PagedResources struct {
	Content      interface{}    `json:"content" binding:"-"`
	PageMetadata pageMetadata   `json:"pageMetadata" binding:"-"`
	Links        []ResourceLink `json:"_links" binding:"-"`
}

type pageMetadata struct {
	TotalElements int      `json:"totalElements"`
	TotalPages    int      `json:"totalPages"`
	Size          int      `json:"size" default:"20"`
	Number        int      `json:"number" default:"0"`
	IndexedBy     string   `json:"indexedBy,omitempty" default:"null"`
	IDs           []string `json:"ids,omitempty"`
}

// ResourceLink defines relationship between resources
type ResourceLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// Resource defines a REST resources links
type Resource struct {
	Links []ResourceLink `json:"_links,omitempty" gorm:"-" binding:"-"`
}

// Resourceable defines a linkable hateoas resource
type Resourceable interface {
	ToResource(baseURL string)
	GetSelfURL(baseURL string) string
}

// ToResources convert receiver page to a hateoas representation
func (page Page) ToResources(baseURL string) PagedResources {
	sortParam := ""
	if page.Pageable.Sort != "" {
		sortParam = "&sort=" + page.Pageable.Sort
	}
	links := []ResourceLink{
		{Href: fmt.Sprintf("%s%s?page=%d&size=%d%s", baseURL, page.BasePath, page.Pageable.Page, page.Pageable.Size, sortParam), Rel: "self"},
	}
	totalPages := int(math.Ceil(float64(page.TotalElements) / float64(page.Pageable.Size)))
	if page.Pageable.Page < totalPages-1 {
		links = append(links, getPageLink(baseURL, page.BasePath, page.Pageable.Page+1, page.Pageable.Size, sortParam))
	}
	if page.Pageable.Page > 0 {
		links = append(links, getPageLink(baseURL, page.BasePath, page.Pageable.Page-1, page.Pageable.Size, sortParam))
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

	reflectValueToResource(reflect.ValueOf(page.Content), baseURL)

	return PagedResources{
		Content:      page.Content,
		PageMetadata: metadata,
		Links:        links,
	}
}

func getPageLink(baseURL, basePath string, page, size int, sortParam string) ResourceLink {
	return ResourceLink{Href: fmt.Sprintf("%s%s?page=%d&size=%d%s", baseURL, basePath, page, size, sortParam), Rel: "next"}
}

func reflectValueToResource(value reflect.Value, baseURL string) {
	if value.Kind() == reflect.Map {
		for _, key := range value.MapKeys() {
			reflectValueToResource(value.MapIndex(key), baseURL)
		}
		return
	}
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			reflectValueToResource(value.Index(i), baseURL)
		}
		return
	}
	if resource, ok := value.Interface().(Resourceable); ok {
		resource.ToResource(baseURL)
	}
}
