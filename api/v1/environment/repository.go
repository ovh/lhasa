package environment

import (
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
)

const (
	defaultPageSize = 20
)

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
	return reflect.TypeOf(v1.Environment{})
}

// GetNewEntityInstance returns a new empty instance of the entity managed by this repository
func (repo *Repository) GetNewEntityInstance() hateoas.Entity {
	return &v1.Environment{}
}

// FindAll returns all entities of the repository type
func (repo *Repository) FindAll() (interface{}, error) {
	return repo.FindBy(map[string]interface{}{})
}

// FindByID gives the details of a particular application
func (repo *Repository) FindByID(id interface{}) (hateoas.Entity, error) {
	env := v1.Environment{}
	err := repo.db.First(&env, id).Error
	if gorm.IsRecordNotFoundError(err) {
		return &env, hateoas.NewEntityDoesNotExistError(env, map[string]interface{}{"id": id})
	}
	if err != nil {
		return nil, err
	}
	return &env, nil
}

// FindOneBySlug fetch a collection of applications matching each criteria
func (repo *Repository) FindOneBySlug(slug string) (*v1.Environment, error) {
	env := v1.Environment{}
	criterias := map[string]interface{}{"slug": slug}
	err := repo.db.First(&env, criterias).Error
	if gorm.IsRecordNotFoundError(err) {
		return &env, hateoas.NewEntityDoesNotExistError(env, criterias)
	}
	if err != nil {
		return nil, err
	}
	return &env, nil
}

// FindBy fetch a collection of applications matching each criteria
func (repo *Repository) FindBy(criterias map[string]interface{}) (interface{}, error) {
	var envs []*v1.Environment
	err := repo.db.Where(criterias).Find(&envs).Error
	return envs, err
}

// FindOneByUnscoped gives the details of a particular environment, even if soft deleted
func (repo *Repository) FindOneByUnscoped(criterias map[string]interface{}) (hateoas.SoftDeletableEntity, error) {
	env := v1.Environment{}
	err := repo.db.Unscoped().Where(criterias).First(&env).Error
	if gorm.IsRecordNotFoundError(err) {
		return &env, hateoas.NewEntityDoesNotExistError(env, criterias)
	}
	return &env, err
}

// FindOneBy fetch the first application matching each criteria
func (repo *Repository) FindOneBy(criterias map[string]interface{}) (hateoas.Entity, error) {
	env := v1.Environment{}
	err := repo.db.Where(criterias).First(&env).Error
	if gorm.IsRecordNotFoundError(err) {
		return &env, hateoas.NewEntityDoesNotExistError(env, criterias)
	}
	return &env, err
}

// Save persists an application to the database
func (repo *Repository) Save(environment hateoas.Entity) error {
	env, err := repo.mustBeEntity(environment)
	if err != nil {
		return err
	}

	if env.ID == 0 {
		return repo.db.Create(env).Error
	}
	return repo.db.Unscoped().Save(env).Error
}

// Remove deletes the application whose GetID is given as a parameter
func (repo *Repository) Remove(env interface{}) error {
	env, err := repo.mustBeEntity(env)
	if err != nil {
		return err
	}

	return repo.db.Delete(env).Error
}

// FindAllPage returns a page of matching entities
func (repo *Repository) FindAllPage(pageable hateoas.Pageable) (hateoas.Page, error) {
	return repo.FindPageBy(pageable, map[string]interface{}{})
}

// FindPageBy returns a page of matching entities
func (repo *Repository) FindPageBy(pageable hateoas.Pageable, criterias map[string]interface{}) (hateoas.Page, error) {
	page := hateoas.NewPage(pageable, defaultPageSize, v1.EnvironmentBasePath)
	var environments []*v1.Environment

	if err := repo.db.Where(criterias).
		Order(page.Pageable.GetGormSortClause()).
		Limit(page.Pageable.Size).
		Offset(page.Pageable.GetOffset()).
		Find(&environments).Error; err != nil {
		return page, err
	}
	page.Content = environments

	count := 0
	if err := repo.db.Model(&v1.Environment{}).Where(criterias).Count(&count).Error; err != nil {
		return page, err
	}
	page.TotalElements = count

	return page, nil
}

// Truncate empties the applications table for testing purposes
func (repo *Repository) Truncate() error {
	return repo.db.Delete(v1.Environment{}).Error
}

func (repo *Repository) mustBeEntity(id interface{}) (*v1.Environment, error) {
	var env *v1.Environment
	switch id := id.(type) {
	case uint:
		env = &v1.Environment{ID: id}
	case *v1.Environment:
		env = id
	case v1.Environment:
		env = &id
	default:
		return nil, hateoas.NewUnsupportedEntityError(env, id)
	}
	return env, nil
}
