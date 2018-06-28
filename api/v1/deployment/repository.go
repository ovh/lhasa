package deployment

import (
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
	"github.com/satori/go.uuid"
)

const (
	defaultPageSize = 20

	queryUndeployByApplicationNameEnvSlug = "UPDATE \"deployments\" AS \"d\" " +
		"SET \"undeployed_at\" = now(), \"updated_at\" = now() " +
		"FROM \"releases\" as \"a\" " +
		"WHERE \"d\".\"deleted_at\" IS NULL " +
		"AND \"a\".\"id\" = \"d\".\"application_id\" " +
		"AND \"a\".\"domain\" = ? " +
		"AND \"a\".\"name\" = ? " +
		"AND \"d\".\"environment_id\" = ? " +
		"AND \"undeployed_at\" IS NULL"
)

// RepositoryFactory defines a repository constructor
type RepositoryFactory func(*gorm.DB) *Repository

// Repository is a repository manager for applications
type Repository struct {
	db *gorm.DB
}

// GetType returns the entity type managed by this repository
func (repo *Repository) GetType() reflect.Type {
	return reflect.TypeOf(v1.Deployment{})
}

// NewRepository creates an application repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// FindPageBy returns a page of matching entities
func (repo *Repository) FindPageBy(pageable hateoas.Pageable, criteria map[string]interface{}) (hateoas.Page, error) {
	page := hateoas.NewPage(pageable, defaultPageSize, v1.DeploymentBasePath)
	var deployments []*v1.Deployment

	// Analyse critarias for extract inline, standard and JSONB ones
	standardCriterias, inlineCriterias, jsonbCriterias := hateoas.CheckFilter(criteria)

	// Apply request
	db := repo.db.Preload("Environment").Preload("Application").Model(v1.Deployment{}).
		Order(page.Pageable.GetGormSortClause()).
		Limit(page.Pageable.Size).
		Offset(page.Pageable.GetOffset())
	db = hateoas.JSONBFilter(db, jsonbCriterias)
	db = hateoas.InlineFilter(db, inlineCriterias)

	if err := db.Find(&deployments, standardCriterias).Error; err != nil {
		return page, err
	}
	page.Content = deployments

	// Build counters
	counter := repo.db.Model(v1.Deployment{}).Where(standardCriterias)
	counter = hateoas.JSONBFilter(counter, jsonbCriterias)
	counter = hateoas.InlineFilter(counter, inlineCriterias)

	count := 0
	if err := counter.Count(&count).Error; err != nil {
		return page, err
	}
	page.TotalElements = count

	return page, nil
}

// Save persists an deployment to the database
func (repo *Repository) Save(deployment hateoas.Entity) error {
	dep, err := repo.mustBeEntity(deployment)
	if err != nil {
		return err
	}

	if dep.ID == 0 {
		publicID, err := uuid.NewV4()
		if err != nil {
			return err
		}
		dep.PublicID = publicID.String()
		return repo.db.Create(dep).Error
	}
	return repo.db.Unscoped().Save(dep).Error
}

// Truncate empties the deployments table for testing purposes
func (repo *Repository) Truncate() error {
	return repo.db.Delete(v1.Deployment{}).Error
}

// Remove deletes the deployment whose GetID is given as a parameter
func (repo *Repository) Remove(dep interface{}) error {
	dep, err := repo.mustBeEntity(dep)
	if err != nil {
		return err
	}

	return repo.db.Delete(dep).Error
}

// FindOneByUnscoped gives the details of a particular deployment, even if soft deleted
func (repo *Repository) FindOneByUnscoped(criteria map[string]interface{}) (hateoas.SoftDeletableEntity, error) {
	dep := &v1.Deployment{}
	err := repo.db.Model(v1.Deployment{}).Unscoped().First(dep, criteria).Error
	if gorm.IsRecordNotFoundError(err) {
		return dep, hateoas.NewEntityDoesNotExistError(dep, criteria)
	}
	return dep, err
}

// FindBy fetch a collection of deployments matching each criteria
func (repo *Repository) FindBy(criteria map[string]interface{}) (interface{}, error) {
	var deps []*v1.Deployment
	err := repo.db.Model(v1.Deployment{}).Find(&deps, criteria).Error
	return deps, err
}

// FindActivesBy fetch a collection of deployments matching each criteria on a given domain and name apps
func (repo *Repository) FindActivesBy(domain string, name string, criterias map[string]interface{}) ([]*v1.Deployment, error) {
	var deps []*v1.Deployment
	err := repo.db.Preload("Environment").Preload("Application").Table("deployments").
		Joins("JOIN releases on releases.id = deployments.application_id").
		Model(v1.Deployment{}).Where("undeployed_at IS NULL AND releases.domain = ? AND releases.name = ?", domain, name).
		Find(&deps, criterias).Error
	return deps, err
}

// FindActivesByVersion fetch a collection of deployments matching each criteria on a given domain, name and version
func (repo *Repository) FindActivesByVersion(domain string, name string, version string, criterias map[string]interface{}) ([]*v1.Deployment, error) {
	var deps []*v1.Deployment
	err := repo.db.Preload("Environment").Preload("Application").Table("deployments").
		Joins("JOIN releases on releases.id = deployments.application_id").
		Model(v1.Deployment{}).Where("undeployed_at IS NULL AND releases.domain = ? AND releases.name = ? AND releases.version = ?", domain, name, version).
		Find(&deps, criterias).Error
	return deps, err
}

// FindOneBy fetch the first deployment matching each criteria
func (repo *Repository) FindOneBy(criteria map[string]interface{}) (hateoas.Entity, error) {
	dep := v1.Deployment{}
	err := repo.db.Where(criteria).First(&dep).Error
	if gorm.IsRecordNotFoundError(err) {
		return &dep, hateoas.NewEntityDoesNotExistError(dep, criteria)
	}
	return &dep, err
}

// UndeployByApplicationEnv updates all deployments attached to a given application regardless version with an undeploy date to now
func (repo *Repository) UndeployByApplicationEnv(domain, name string, envID uint) error {
	return repo.db.Exec(queryUndeployByApplicationNameEnvSlug, domain, name, envID).Error
}

func (repo *Repository) mustBeEntity(id interface{}) (*v1.Deployment, error) {
	var dep *v1.Deployment
	switch id := id.(type) {
	case uint:
		dep = &v1.Deployment{ID: id}
	case *v1.Deployment:
		dep = id
	case v1.Deployment:
		dep = &id
	default:
		return nil, hateoas.NewUnsupportedEntityError(dep, id)
	}
	return dep, nil
}
