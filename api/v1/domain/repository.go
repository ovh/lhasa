package domain

import (
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
)

const (
	defaultPageSize = 20
)

// Repository is a repository manager for domains
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
	return reflect.TypeOf(v1.Domain{})
}

// GetNewEntityInstance returns a new empty instance of the entity managed by this repository
func (repo *Repository) GetNewEntityInstance() hateoas.Entity {
	return &v1.Domain{}
}

// FindAll returns all entities of the repository type
func (repo *Repository) FindAll() (interface{}, error) {
	return repo.FindBy(map[string]interface{}{})
}

// FindAllPage returns a page of matching entities
func (repo *Repository) FindAllPage(pageable hateoas.Pageable) (hateoas.Page, error) {
	return repo.FindPageBy(pageable, map[string]interface{}{})
}

// FindPageBy returns a page of matching entities
func (repo *Repository) FindPageBy(pageable hateoas.Pageable, criterias map[string]interface{}) (hateoas.Page, error) {
	page := hateoas.NewPage(pageable, defaultPageSize, v1.DomainBasePath)
	var domainNames []string

	if err := repo.db.Model(&v1.ApplicationVersion{}).
		Where(criterias).
		Order(page.Pageable.GetSortClause()).
		Limit(page.Pageable.Size).
		Offset(page.Pageable.GetOffset()).
		Pluck("DISTINCT \"applications\".\"domain\" as \"name\"", &domainNames).Error; err != nil {
		return page, err
	}

	var domains []*v1.Domain
	for _, value := range domainNames {
		domains = append(domains, &v1.Domain{Name: value})
	}
	page.Content = domains

	count := 0
	rows, err := repo.db.Raw("select count(distinct \"applications\".\"domain\") \"totalElements\" from \"applications\" where \"applications\".\"deleted_at\" is null").Rows()
	if err != nil {
		return page, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&count)
	page.TotalElements = count

	return page, nil
}

// Save persists an domain to the database
func (repo *Repository) Save(_ hateoas.Entity) error {
	return errors.NotSupportedf("operation not supported")
}

// Truncate empties the domains table for testing purposes
func (repo *Repository) Truncate() error {
	return errors.NotSupportedf("operation not supported")
}

// Remove deletes the domain whose GetID is given as a parameter
func (repo *Repository) Remove(_ interface{}) error {
	return errors.NotSupportedf("operation not supported")
}

// FindByID gives the details of a particular domain
func (repo *Repository) FindByID(id interface{}) (hateoas.Entity, error) {
	return nil, errors.NotSupportedf("operation not supported")
}

// FindBy fetch a collection of domains matching each criteria
func (repo *Repository) FindBy(criterias map[string]interface{}) (interface{}, error) {
	return nil, errors.NotSupportedf("operation not supported")
}

// FindOneBy find by criterias
func (repo *Repository) FindOneBy(criterias map[string]interface{}) (hateoas.Entity, error) {
	var app v1.ApplicationVersion
	err := repo.db.First(&app, criterias).Error
	if gorm.IsRecordNotFoundError(err) {
		return &app, hateoas.NewEntityDoesNotExistError(app, criterias)
	}
	if err != nil {
		return nil, err
	}
	return &v1.Domain{Name: app.Domain}, nil
}
