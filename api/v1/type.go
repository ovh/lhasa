package v1

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"github.com/ovh/lhasa/api/graphapi"
	"github.com/ovh/lhasa/api/hateoas"
)

// ApplicationBasePath is the URL base path for this resource
const ApplicationBasePath = "/applications"

// DeploymentBasePath is the URL base path for this resource
const DeploymentBasePath = "/deployments"

// EnvironmentBasePath is the URL base path for this resource
const EnvironmentBasePath = "/environments"

// DomainBasePath is the URL base path for this resource
const DomainBasePath = "/domains"

// ContentBasePath is the URL base path for this resource
const ContentBasePath = "/content"

// MediaResource defines a media resource behaviour
type MediaResource interface {
	GetContentType() string
	GetBytes() []byte
	SetBytes([]byte)
}

// Content define a content resource
type Content struct {
	ID          uint       `json:"-" binding:"-" gorm:"auto increment"`
	Name        string     `json:"-" path:"name" description:"Content data"`
	ContentType string     `json:"-" header:"content-type" description:"Content mime type"`
	Locale      string     `json:"-" path:"locale" description:"Locale"`
	Body        []byte     `json:"-" description:"Content body"`
	CreatedAt   time.Time  `json:"-" binding:"-"`
	UpdatedAt   time.Time  `json:"-" binding:"-"`
	DeletedAt   *time.Time `json:"-" binding:"-"`
	hateoas.Resource
}

// GetContentType get content type
func (p *Content) GetContentType() string {
	return p.ContentType
}

// GetBytes get all bytes
func (p *Content) GetBytes() []byte {
	return p.Body
}

// SetBytes set all bytes
func (p *Content) SetBytes(data []byte) {
	p.Body = data
}

// Domain define a business domain
type Domain struct {
	Name string `json:"name" path:"domain" description:"Application Domain"`
	hateoas.Resource
}

// Application defines an application
type Application struct {
	Domain string `path:"domain" description:"Application Domain"`
	Name   string `path:"name" description:"Application Name"`
}

// TableName Set ApplicationVersion's table name to be `applications`
func (ApplicationVersion) TableName() string {
	return "applications"
}

// ApplicationVersion defines the model properties of an application version
type ApplicationVersion struct {
	ID           uint            `json:"-" gorm:"auto increment" binding:"-"`
	Domain       string          `json:"domain" gorm:"not null;type:varchar(255);unique_index:idx_applications_domain_name_version;default:''" path:"domain" description:"Application Domain"`
	Name         string          `json:"name" gorm:"not null;type:varchar(255);unique_index:idx_applications_domain_name_version;default:''" path:"name" description:"Application Name"`
	Version      string          `json:"version" gorm:"not null;type:varchar(255);unique_index:idx_applications_domain_name_version;default:''" path:"version" description:"Application Version"`
	Manifest     *postgres.Jsonb `json:"manifest"`
	Tags         pq.StringArray  `json:"tags,omitempty" gorm:"type:varchar(255)[]"`
	Dependencies []Dependency    `json:"-"`
	Deployments  []Deployment    `json:"-"`
	CreatedAt    time.Time       `json:"_createdAt" binding:"-"`
	UpdatedAt    time.Time       `json:"_updatedAt" binding:"-"`
	DeletedAt    *time.Time      `json:"-" binding:"-"`
	hateoas.Resource
}

// Dependency defines a inter-application link
type Dependency struct {
	ID       uint `json:"-" gorm:"auto increment"`
	Owner    ApplicationVersion
	OwnerID  uint `json:"-" gorm:"type:bigint;not null;default:0"`
	Target   ApplicationVersion
	TargetID uint `json:"-" gorm:"type:bigint;not null;default:0"`
}

// Deployment is an application version instance on a given environment
type Deployment struct {
	ID            uint                `json:"-" gorm:"auto increment" binding:"-"`
	PublicID      string              `json:"id" path:"public_id" gorm:"type:varchar(255);not null;unique" validate:"omitempty,uuid4" binding:"omitempty,uuid4" description:"Deployment public identifier"`
	ApplicationID uint                `json:"-" gorm:"not null;type:bigint;default:0"`
	Application   *ApplicationVersion `json:"-"`
	EnvironmentID uint                `json:"-" gorm:"not null;type:bigint;default:0"`
	Environment   *Environment        `json:"-"`
	Dependencies  postgres.Jsonb      `json:"dependencies,omitempty" binding:"-"`
	Properties    postgres.Jsonb      `json:"properties,omitempty"`
	UndeployedAt  *time.Time          `json:"undeployedAt,omitempty" binding:"-"`
	CreatedAt     time.Time           `json:"_createdAt" binding:"-"`
	UpdatedAt     time.Time           `json:"_updatedAt" binding:"-"`
	DeletedAt     *time.Time          `json:"-" binding:"-"`
	hateoas.Resource
}

// DeploymentDependency defines a inter-deployment link
type DeploymentDependency struct {
	TargetID string `json:"target"`
	Type     string `json:"type"`
}

// Environment is a target deployment environment
type Environment struct {
	ID         uint           `json:"-" gorm:"auto increment" binding:"-"`
	Slug       string         `json:"slug" gorm:"type:varchar(255);unique;not null;default:''" path:"slug" description:"Environment identifier"`
	Name       string         `json:"name" gorm:"type:varchar(255)"`
	Properties postgres.Jsonb `json:"properties,omitempty"`
	CreatedAt  time.Time      `json:"_createdAt" binding:"-"`
	UpdatedAt  time.Time      `json:"_updatedAt" binding:"-"`
	DeletedAt  *time.Time     `json:"-" binding:"-"`
	hateoas.Resource
}

// ApplicationVersionGraph is an application version instance on a given environment
type ApplicationVersionGraph struct {
	PublicID  string          `json:"id" description:"Application public identifier"`
	Domain    string          `json:"domain" description:"Application Domain"`
	Name      string          `json:"name" description:"Application Name"`
	Version   string          `json:"version" description:"Application Version"`
	Manifest  json.RawMessage `json:"manifest"`
	Tags      []string        `json:"tags,omitempty" `
	CreatedAt time.Time       `json:"_createdAt" `
	UpdatedAt time.Time       `json:"_updatedAt" `
}

// GetID return ID
func (a *ApplicationVersionGraph) GetID() string {
	return a.PublicID
}

// GetName get all dependencies
func (a *ApplicationVersionGraph) GetName() string {
	return a.Domain + "/" + a.Name + " v" + a.Version
}

// GetType get entity type
func (a *ApplicationVersionGraph) GetType() string {
	return "application"
}

// GetEdges get all dependencies
func (a *ApplicationVersionGraph) GetEdges() []graphapi.ImplementEdge {
	return []graphapi.ImplementEdge{}
}

// DeploymentGraph is an application version instance on a given environment
type DeploymentGraph struct {
	PublicID     string                 `json:"id" description:"Deployment public identifier"`
	Application  graphapi.ImplementNode `json:"application"`
	Environment  graphapi.ImplementNode `json:"environnement"`
	Dependencies json.RawMessage        `json:"dependencies,omitempty" binding:"-"`
	Properties   json.RawMessage        `json:"properties,omitempty"`
	UndeployedAt *time.Time             `json:"undeployedAt,omitempty" binding:"-"`
	CreatedAt    time.Time              `json:"_createdAt" binding:"-"`
	UpdatedAt    time.Time              `json:"_updatedAt" binding:"-"`
}

// GetID get deployment ID
func (d *DeploymentGraph) GetID() string {
	return d.PublicID
}

// GetName get all dependencies
func (d *DeploymentGraph) GetName() string {
	return d.Application.GetName() + " on (" + d.Environment.GetName() + ")"
}

// GetType get entity type
func (d *DeploymentGraph) GetType() string {
	return "deployment"
}

// GetEdges get all dependencies
func (d *DeploymentGraph) GetEdges() []graphapi.ImplementEdge {
	dependencies := []graphapi.ImplementEdge{}
	if d.Dependencies != nil {
		rawDependencies := []DeploymentDependency{}
		json.Unmarshal(d.Dependencies, &rawDependencies)
		for _, entity := range rawDependencies {
			uid, err := uuid.NewV4()
			if err == nil {
				edge := graphapi.Edge{
					ID:   uid.String(),
					From: d.PublicID,
					To:   entity.TargetID,
					Type: entity.Type,
					Properties: map[string]interface{}{
						"Type": entity.Type,
					},
				}
				dependencies = append(dependencies, &edge)
			}
		}
	}
	return dependencies
}

// EnvironmentGraph is an application version instance on a given environment
type EnvironmentGraph struct {
	PublicID   string          `json:"id" description:"Environnement public identifier"`
	Slug       string          `json:"slug" description:"Environment identifier"`
	Name       string          `json:"name"`
	Properties json.RawMessage `json:"properties,omitempty"`
	CreatedAt  time.Time       `json:"_createdAt" binding:"-"`
	UpdatedAt  time.Time       `json:"_updatedAt" binding:"-"`
}

// GetID return ID
func (e *EnvironmentGraph) GetID() string {
	return e.PublicID
}

// GetName get all dependencies
func (e *EnvironmentGraph) GetName() string {
	return e.Slug + "/" + e.Name
}

// GetType get entity type
func (e *EnvironmentGraph) GetType() string {
	return "environnement"
}

// GetEdges get all dependencies
func (e *EnvironmentGraph) GetEdges() []graphapi.ImplementEdge {
	return []graphapi.ImplementEdge{}
}
