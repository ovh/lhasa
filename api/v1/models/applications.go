package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	"github.com/ovh/lhasa/api/repositories"
)

// ApplicationBasePath is the URL base path for this resource
const ApplicationBasePath = "/applications"

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
	repositories.Resource
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
	app.Resource.Links = []repositories.ResourceLink{
		{Rel: "self", Href: app.GetSelfURL(baseURL)},
		{Rel: "deployments", Href: app.GetSelfURL(baseURL) + "/deployments"},
		{Rel: "deploy", Href: app.GetSelfURL(baseURL) + "/deploy/:environment"}}
}

// GetSelfURL implements Resourceable
func (app *Application) GetSelfURL(baseURL string) string {
	return fmt.Sprintf("%s%s/%s/%s/%s", baseURL, ApplicationBasePath, app.Domain, app.Name, app.Version)
}
