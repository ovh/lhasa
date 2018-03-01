package models

import (
	"time"
)

// Application defines the model properties of an application
type Application struct {
	ID           uint             `json:"-" gorm:"auto increment"`
	Domain       string           `json:"profile" gorm:"not null;unique_index:idx_unique_fullname;default:''" binding:"required"`
	Name         string           `json:"name" gorm:"not null;unique_index:idx_unique_fullname;default:''" binding:"required"`
	Version      string           `json:"version" gorm:"not null;unique_index:idx_unique_fullname;default:''"`
	Type         string           `json:"type,omitempty"`
	Artefact     string           `json:"artefact,omitempty"`
	Language     string           `json:"language,omitempty"`
	Description  string           `json:"description,omitempty"`
	Repository   string           `json:"repository,omitempty"`
	Package      string           `json:"package,omitempty"`
	Authors      []PersonInfo     `json:"authors,omitempty"`
	Support      SupportInfo      `json:"support,omitempty" gorm:"embedded;embedded_prefix:support_"`
	Requires     []DependencyInfo `json:"requires,omitempty"`
	EventBus     EventBusInfo     `json:"eventbus,omitempty" gorm:"embedded;embedded_prefix:eventbus_"`
	VaultAliases string           `json:"vault,omitempty"`
	Tags         string           `json:"tags,omitempty"`
	CreatedAt    time.Time        `json:"-"`
	UpdatedAt    time.Time        `json:"-"`
	DeletedAt    *time.Time       `json:"-"`
}

// PersonInfo defines a person and her/his role in the project
type PersonInfo struct {
	ID      uint   `json:"-"`
	Name    string `json:"name,omitempty"`
	Country string `json:"country,omitempty"`
	Email   string `json:"email,omitempty"`
	Role    string `json:"role,omitempty"`
}

// SupportInfo defines support informations
type SupportInfo struct {
	Email  string `json:"email,omitempty"`
	Issues string `json:"issues,omitempty"`
	Docs   string `json:"docs,omitempty"`
}

// DependencyInfo defines project requirements
type DependencyInfo struct {
	ID       uint   `json:"-"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Version  string `json:"version,omitempty"`
	Critical bool   `json:"critical,omitempty"`
}

// EventBusInfo defines produced and consumed events
type EventBusInfo struct {
	Produces string `json:"produces,omitempty"`
	Consumes string `json:"consumes,omitempty"`
}
