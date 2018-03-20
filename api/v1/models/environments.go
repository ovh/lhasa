package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ovh/lhasa/api/repositories"
)

// EnvironmentBasePath is the URL base path for this resource
const EnvironmentBasePath = "/environments"

// Environment is a target deployment environment
type Environment struct {
	ID         uint           `json:"-" gorm:"auto increment"`
	Slug       string         `json:"slug" gorm:"type:varchar(255);unique;not null;default:''" path:"slug"`
	Name       string         `json:"name" gorm:"type:varchar(255)"`
	Properties postgres.Jsonb `json:"properties,omitempty"`
	CreatedAt  time.Time      `json:"_createdAt"`
	UpdatedAt  time.Time      `json:"_updatedAt"`
	DeletedAt  *time.Time     `json:"-"`
	repositories.Resource
}

// GetID returns the public ID of the entity
func (env *Environment) GetID() string {
	return env.Slug
}

// SetID sets up the new ID of the entity
func (env *Environment) SetID(id string) error {
	env.Slug = id
	return nil
}

// GetDeletedAt implements SoftDeletableEntity
func (env *Environment) GetDeletedAt() *time.Time {
	return env.DeletedAt
}

// ToResource implements Resourceable
func (env *Environment) ToResource(baseURL string) {
	env.Resource.Links = []repositories.ResourceLink{{Rel: "self", Href: env.GetSelfURL(baseURL)}}
}

// GetSelfURL implements Resourceable
func (env *Environment) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, EnvironmentBasePath, env.Slug)
}
