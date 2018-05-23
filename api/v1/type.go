package v1

import (
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
