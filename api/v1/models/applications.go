package models

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

// Application defines the model properties of an application
type Application struct {
	ID           uint            `json:"-" gorm:"auto increment"`
	Domain       string          `json:"domain" gorm:"not null;type:varchar(255);unique_index:idx_applications_domain_name_version;default:''" path:"domain"`
	Name         string          `json:"name" gorm:"not null;type:varchar(255);unique_index:idx_applications_domain_name_version;default:''" path:"name"`
	Version      string          `json:"version" gorm:"not null;type:varchar(255);unique_index:idx_applications_domain_name_version;default:''" path:"version"`
	Manifest     *postgres.Jsonb `json:"manifest"`
	Tags         pq.StringArray  `json:"-" gorm:"type:varchar(255)[]"`
	Dependencies []Dependency    `json:"-" gorm:"foreignkey:OwnerID"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	DeletedAt    *time.Time      `json:"-"`
}

// Dependency defines a inter-application link
type Dependency struct {
	ID       uint `json:"-" gorm:"auto increment"`
	Owner    Application
	OwnerID  uint `json:"-" gorm:"type:serial"`
	Target   Application
	TargetID uint `json:"-" gorm:"type:serial"`
}
