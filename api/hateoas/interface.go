package hateoas

import (
	"fmt"
	"math"
	"reflect"
	"time"
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

// Repository defines a normal repository
type Repository interface {
	FindAll() (interface{}, error)
	FindByID(interface{}) (Entity, error)
	FindBy(map[string]interface{}) (interface{}, error)
	FindOneBy(map[string]interface{}) (Entity, error)
	Save(Entity) error
	Remove(interface{}) error
	GetNewEntityInstance() Entity
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

// SoftDeletableRepository defines a repository that handles soft delete
type SoftDeletableRepository interface {
	Repository
	FindOneByUnscoped(criterias map[string]interface{}) (SoftDeletableEntity, error)
}

// Pageable defines pagination criterias
type Pageable struct {
	Page      int    `json:"page" form:"page"`
	Size      int    `json:"size" form:"size"`
	Sort      string `json:"sort" form:"sort"`
	Query     string `json:"q" form:"q"`
	IndexedBy string `json:"indexedBy" form:"indexedBy"`
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
	Content      interface{}    `json:"content"`
	PageMetadata pageMetadata   `json:"pageMetadata"`
	Links        []ResourceLink `json:"_links"`
}

type pageMetadata struct {
	TotalElements int      `json:"totalElements"`
	TotalPages    int      `json:"totalPages"`
	Size          int      `json:"size"`
	Number        int      `json:"number"`
	IndexedBy     string   `json:"indexedBy,omitempty"`
	IDs           []string `json:"ids,omitempty"`
}

// ResourceLink defines relationship between resources
type ResourceLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// Resource defines a REST resources links
type Resource struct {
	Links []ResourceLink `json:"_links,omitempty" gorm:"-"`
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
		links = append(links, ResourceLink{Href: fmt.Sprintf("%s%s?page=%d&size=%d%s", baseURL, page.BasePath, page.Pageable.Page+1, page.Pageable.Size, sortParam), Rel: "next"})
	}
	if page.Pageable.Page > 0 {
		links = append(links, ResourceLink{Href: fmt.Sprintf("%s%s?page=%d&size=%d%s", baseURL, page.BasePath, page.Pageable.Page-1, page.Pageable.Size, sortParam), Rel: "prev"})
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

func reflectValueToResource(value reflect.Value, baseURL string) {
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Map {
		resource, ok := value.Interface().(Resourceable)
		if ok {
			resource.ToResource(baseURL)
		}
		return
	}
	for i := 0; i < value.Len(); i++ {
		reflectValueToResource(value.Index(i), baseURL)
	}
}
