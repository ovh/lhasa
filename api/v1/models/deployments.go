package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/ovh/lhasa/api/repositories"
)

// DeploymentBasePath is the URL base path for this resource
const DeploymentBasePath = "/deployments"

// Deployment is an application version instance on a given environment
type Deployment struct {
	ID            uint           `json:"-" gorm:"auto increment"`
	PublicID      string         `json:"id" gorm:"type:varchar(255);not null;unique"`
	ApplicationID uint           `json:"-" gorm:"not null;type:bigint;default:0"`
	Application   Application    `json:"-"`
	EnvironmentID uint           `json:"-" gorm:"not null;type:bigint;default:0"`
	Environment   Environment    `json:"-"`
	Properties    postgres.Jsonb `json:"properties,omitempty"`
	UndeployedAt  *time.Time     `json:"undeployedAt"`
	CreatedAt     time.Time      `json:"_createdAt"`
	UpdatedAt     time.Time      `json:"_updatedAt"`
	DeletedAt     *time.Time     `json:"-"`
	repositories.Resource
}

// GetID returns the public ID of the entity
func (dep *Deployment) GetID() string {
	return dep.PublicID
}

// SetID sets up the new ID of the entity
func (dep *Deployment) SetID(id string) error {
	dep.PublicID = id
	return nil
}

// GetDeletedAt implements SoftDeletableEntity
func (dep *Deployment) GetDeletedAt() *time.Time {
	return dep.DeletedAt
}

// ToResource implements Resourceable
func (dep *Deployment) ToResource(baseURL string) {
	dep.Resource.Links = []repositories.ResourceLink{
		{Rel: "self", Href: dep.GetSelfURL(baseURL)},
		{Rel: "environment", Href: dep.Environment.GetSelfURL(baseURL)},
		{Rel: "application", Href: dep.Application.GetSelfURL(baseURL)},
	}
}

// GetSelfURL implements Resourceable
func (dep *Deployment) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, DeploymentBasePath, dep.PublicID)
}
