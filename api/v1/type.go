package v1

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	"github.com/ovh/lhasa/api/hateoas"
)

// ApplicationBasePath is the URL base path for this resource
const ApplicationBasePath = "/applications"

// DeploymentBasePath is the URL base path for this resource
const DeploymentBasePath = "/deployments"

// EnvironmentBasePath is the URL base path for this resource
const EnvironmentBasePath = "/environments"

// Application defines the model properties of an application
type Application struct {
	ID           uint            `json:"-" gorm:"auto increment"`
	Domain       string          `json:"domain" gorm:"not null;type:varchar(255);unique_index:idx_applications_domain_name_version;default:''" path:"domain"`
	Name         string          `json:"name" gorm:"not null;type:varchar(255);unique_index:idx_applications_domain_name_version;default:''" path:"name"`
	Version      string          `json:"version" gorm:"not null;type:varchar(255);unique_index:idx_applications_domain_name_version;default:''" path:"version"`
	Manifest     *postgres.Jsonb `json:"manifest"`
	Tags         pq.StringArray  `json:"-" gorm:"type:varchar(255)[]"`
	Dependencies []Dependency    `json:"-" gorm:"foreignkey:OwnerID"`
	Deployments  []Deployment    `json:"-"`
	CreatedAt    time.Time       `json:"_createdAt"`
	UpdatedAt    time.Time       `json:"_updatedAt"`
	DeletedAt    *time.Time      `json:"-"`
	hateoas.Resource
}

// Dependency defines a inter-application link
type Dependency struct {
	ID       uint `json:"-" gorm:"auto increment"`
	Owner    Application
	OwnerID  uint `json:"-" gorm:"type:bigint;not null;default:0"`
	Target   Application
	TargetID uint `json:"-" gorm:"type:bigint;not null;default:0"`
}

// GetID returns the public ID of the entity
func (app *Application) GetID() string {
	return string(app.ID)
}

// SetID sets up the new ID of the entity
func (app *Application) SetID(id string) error {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	app.ID = uint(idInt)
	return nil
}

// GetDeletedAt implements SoftDeletableEntity
func (app *Application) GetDeletedAt() *time.Time {
	return app.DeletedAt
}

// ToResource implements Resourceable
func (app *Application) ToResource(baseURL string) {
	app.Resource.Links = []hateoas.ResourceLink{
		{Rel: "self", Href: app.GetSelfURL(baseURL)},
		{Rel: "deployments", Href: app.GetSelfURL(baseURL) + "/deployments"},
		{Rel: "deploy", Href: app.GetSelfURL(baseURL) + "/deploy/:environment"}}
}

// GetSelfURL implements Resourceable
func (app *Application) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s/%s/%s", baseURL, ApplicationBasePath, app.Domain, app.Name, app.Version)
}

// Deployment is an application version instance on a given environment
type Deployment struct {
	ID            uint           `json:"-" gorm:"auto increment"`
	PublicID      string         `json:"id" gorm:"type:varchar(255);not null;unique"`
	ApplicationID uint           `json:"-" gorm:"not null;type:bigint;default:0"`
	Application   *Application   `json:"-"`
	EnvironmentID uint           `json:"-" gorm:"not null;type:bigint;default:0"`
	Environment   *Environment   `json:"-"`
	Properties    postgres.Jsonb `json:"properties,omitempty"`
	UndeployedAt  *time.Time     `json:"undeployedAt,omitempty"`
	CreatedAt     time.Time      `json:"_createdAt"`
	UpdatedAt     time.Time      `json:"_updatedAt"`
	DeletedAt     *time.Time     `json:"-"`
	hateoas.Resource
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
	dep.Resource.Links = []hateoas.ResourceLink{{Rel: "self", Href: dep.GetSelfURL(baseURL)}}
	if dep.Environment != nil {
		dep.Resource.Links = append(dep.Resource.Links, hateoas.ResourceLink{Rel: "environment", Href: dep.Environment.GetSelfURL(baseURL)})
	}
	if dep.Application != nil {
		dep.Resource.Links = append(dep.Resource.Links, hateoas.ResourceLink{Rel: "application", Href: dep.Application.GetSelfURL(baseURL)})
	}
}

// GetSelfURL implements Resourceable
func (dep *Deployment) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, DeploymentBasePath, dep.PublicID)
}

// Environment is a target deployment environment
type Environment struct {
	ID         uint           `json:"-" gorm:"auto increment"`
	Slug       string         `json:"slug" gorm:"type:varchar(255);unique;not null;default:''" path:"slug"`
	Name       string         `json:"name" gorm:"type:varchar(255)"`
	Properties postgres.Jsonb `json:"properties,omitempty"`
	CreatedAt  time.Time      `json:"_createdAt"`
	UpdatedAt  time.Time      `json:"_updatedAt"`
	DeletedAt  *time.Time     `json:"-"`
	hateoas.Resource
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
	env.Resource.Links = []hateoas.ResourceLink{{Rel: "self", Href: env.GetSelfURL(baseURL)}}
}

// GetSelfURL implements Resourceable
func (env *Environment) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s", baseURL, EnvironmentBasePath, env.Slug)
}
