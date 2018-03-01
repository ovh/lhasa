package models

import (
	"github.com/ovh/lhasa/api/db"
)

// Application defines the model properties of an application
type Application struct {
	ID            int    `json:"id" gorm:"primary_key;auto_increment"`
	Domain        string `json:"domain" validate:"presence,min=1,max=32" gorm:"not null,unique_index:idx_unique_fullname"`
	Name          string `json:"name" validate:"presence,min=1,max=32" gorm:"not null,unique_index:idx_unique_fullname"`
	Type          string `json:"type" validate:"presence,min=1,max=32"` // possible values: service, job, cronjob, ...
	Language      string `json:"language"`                              // e.g. python, golang, perl, js, php, java
	RepositoryURL string `json:"repositoryurl" validate:"url"`
	AvatarURL     string `json:"avatarurl" validate:"url"`
	Description   string `json:"description" sql:"type:text;"`
}

// MigrateApplications runs automated gorm migrations
func MigrateApplications() error {
	return db.DB().AutoMigrate(&Application{}).Error
}

// ListApplications provides a paginated list of potentially filtered applications
func ListApplications(query string, start int, size int) ([]Application, error) {
	var applications []Application
	err := db.DB().Offset(start).Limit(size).Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&applications).Error
	return applications, err
}

// CreateApplication adds an application to the database
func CreateApplication(app *Application) error {
	return db.DB().Create(&app).Error
}

// FlushApplications empties the applications table for testing purposes
func FlushApplications() error {
	return db.DB().Exec("delete from Applications").Error
}

// DeleteApplication deletes the application whose ID is given as a parameter
func DeleteApplication(id int) (*Application, error) {
	app := &Application{ID: id}
	if err := db.DB().Delete(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

// DetailApplication gives the details of a particular application
func DetailApplication(id int) (*Application, error) {
	app := &Application{}
	if err := db.DB().First(app, id).Error; err != nil {
		return nil, err
	}
	return app, nil
}
