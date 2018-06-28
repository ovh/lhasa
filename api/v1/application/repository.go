package application

import (
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
)

const (
	defaultPageSize = 20
)

// FindOneByUniqueKey declare a method allowing to find an application by domain, name and version
type FindOneByUniqueKey interface {
	FindOneByDomainNameVersion(string, string, string) (*v1.Release, error)
}

// RepositoryFactory defines a repository constructor
type RepositoryFactory func(*gorm.DB) *Repository

// Repository is a repository manager for applications
type Repository struct {
	db *gorm.DB
}

// NewRepository creates an application repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetType returns the entity type managed by this repository
func (repo *Repository) GetType() reflect.Type {
	return reflect.TypeOf(v1.Release{})
}

// FindPageBy returns a page of matching entities
func (repo *Repository) FindPageBy(pageable hateoas.Pageable, criteria map[string]interface{}) (hateoas.Page, error) {
	page := hateoas.NewPage(pageable, defaultPageSize, v1.ApplicationBasePath)
	var applications []*v1.Release

	if err := repo.db.
		Where(criteria).
		Order(page.Pageable.GetGormSortClause()).
		Limit(page.Pageable.Size).
		Offset(page.Pageable.GetOffset()).
		Find(&applications).Error; err != nil {
		return page, err
	}
	page.Content = applications

	count := 0
	if err := repo.db.Model(&v1.Release{}).Where(criteria).Count(&count).Error; err != nil {
		return page, err
	}
	page.TotalElements = count

	if pageable.IndexedBy != "" {
		currentIndex := map[string][]*v1.Release{}
		ids := map[string]bool{}
		for _, application := range applications {
			indexedField, err := repo.getIndexedField(pageable.IndexedBy, application)
			if err != nil {
				return page, err
			}
			currentIndex[indexedField] = append(currentIndex[indexedField], application)
			ids[indexedField] = true
		}
		page.Content = currentIndex
		for id := range ids {
			page.IDs = append(page.IDs, id)
		}
	}

	return page, nil
}

func (repo *Repository) getIndexedField(field string, application *v1.Release) (string, error) {
	switch field {
	case "version":
		return application.Version, nil
	case "profile":
	case "domain":
		return application.Domain, nil
	}
	return "", hateoas.NewUnsupportedIndexError(field, "version", "domain")
}

// Save persists an application to the database
func (repo *Repository) Save(application hateoas.Entity) error {
	app, err := repo.mustBeEntity(application)
	if err != nil {
		return err
	}

	if app.ID == 0 {
		return repo.db.Create(app).Error
	}
	return repo.db.Unscoped().Save(app).Error
}

// Truncate empties the applications table for testing purposes
func (repo *Repository) Truncate() error {
	return repo.db.Delete(v1.Release{}).Error
}

// Remove deletes the application whose GetID is given as a parameter
func (repo *Repository) Remove(app interface{}) error {
	app, err := repo.mustBeEntity(app)
	if err != nil {
		return err
	}

	return repo.db.Delete(app).Error
}

// FindOneByUnscoped gives the details of a particular application, even if soft deleted
func (repo *Repository) FindOneByUnscoped(criteria map[string]interface{}) (hateoas.SoftDeletableEntity, error) {
	app := v1.Release{}
	err := repo.db.Unscoped().Where(criteria).First(&app).Error
	if gorm.IsRecordNotFoundError(err) {
		return &app, hateoas.NewEntityDoesNotExistError(app, criteria)
	}
	return &app, err
}

// FindBy fetch a collection of applications matching each criteria
func (repo *Repository) FindBy(criteria map[string]interface{}) (interface{}, error) {
	var apps []*v1.Release
	err := repo.db.Where(criteria).Find(&apps).Error
	return apps, err
}

// FindOneByDomainNameVersion fetch the first application matching each criteria
func (repo *Repository) FindOneByDomainNameVersion(domain, name, version string) (*v1.Release, error) {
	app := v1.Release{}
	criteria := map[string]interface{}{
		"domain":  domain,
		"name":    name,
		"version": version,
	}
	err := repo.db.First(&app, criteria).Error
	if gorm.IsRecordNotFoundError(err) {
		return &app, hateoas.NewEntityDoesNotExistError(app, criteria)
	}
	return &app, err
}

// FindOneByDomainName fetch the first application matching each criteria
func (repo *Repository) FindOneByDomainName(domain, name string) (*v1.Release, error) {
	app := v1.Release{}
	criteria := map[string]interface{}{
		"domain": domain,
		"name":   name,
	}
	err := repo.db.First(&app, criteria).Error
	if gorm.IsRecordNotFoundError(err) {
		return &app, hateoas.NewEntityDoesNotExistError(app, criteria)
	}
	return &app, err
}

// FindOneBy find by criteria
func (repo *Repository) FindOneBy(criteria map[string]interface{}) (hateoas.Entity, error) {
	app := v1.Release{}
	err := repo.db.Where(criteria).First(&app).Error
	if gorm.IsRecordNotFoundError(err) {
		return &app, hateoas.NewEntityDoesNotExistError(app, criteria)
	}
	return &app, err
}

func (repo *Repository) mustBeEntity(id interface{}) (*v1.Release, error) {
	var app *v1.Release
	switch id := id.(type) {
	case uint:
		app = &v1.Release{ID: id}
	case *v1.Release:
		app = id
	case v1.Release:
		app = &id
	default:
		return nil, hateoas.NewUnsupportedEntityError(app, id)
	}
	return app, nil
}
