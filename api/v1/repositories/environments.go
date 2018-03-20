package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/ovh/lhasa/api/repositories"
	"github.com/ovh/lhasa/api/v1/models"
)

const (
	environmentsDefaultPageSize = 20
)

// EnvironmentRepository is a repository manager for applications
type EnvironmentRepository struct {
	db *gorm.DB
}

// NewEnvironmentRepository creates an application repository
func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{
		db: db,
	}
}

// GetNewEntityInstance returns a new empty instance of the entity managed by this repository
func (repo *EnvironmentRepository) GetNewEntityInstance() repositories.Entity {
	return &models.Environment{}
}

// FindAll returns all entities of the repository type
func (repo *EnvironmentRepository) FindAll() (interface{}, error) {
	return repo.FindBy(map[string]interface{}{})
}

// FindByID gives the details of a particular application
func (repo *EnvironmentRepository) FindByID(id interface{}) (repositories.Entity, error) {
	env := models.Environment{}
	err := repo.db.First(&env, id).Error
	if gorm.IsRecordNotFoundError(err) {
		return &env, repositories.NewEntityDoesNotExistError(env, map[string]interface{}{"id": id})
	}
	if err != nil {
		return nil, err
	}
	return &env, nil
}

// FindOneBySlug fetch a collection of applications matching each criteria
func (repo *EnvironmentRepository) FindOneBySlug(slug string) (*models.Environment, error) {
	env := models.Environment{}
	criterias := map[string]interface{}{"slug": slug}
	err := repo.db.First(&env, criterias).Error
	if gorm.IsRecordNotFoundError(err) {
		return &env, repositories.NewEntityDoesNotExistError(env, criterias)
	}
	if err != nil {
		return nil, err
	}
	return &env, nil
}

// FindBy fetch a collection of applications matching each criteria
func (repo *EnvironmentRepository) FindBy(criterias map[string]interface{}) (interface{}, error) {
	var envs []*models.Environment
	err := repo.db.Where(criterias).Find(&envs).Error
	return envs, err
}

// FindOneByUnscoped gives the details of a particular environment, even if soft deleted
func (repo *EnvironmentRepository) FindOneByUnscoped(criterias map[string]interface{}) (repositories.SoftDeletableEntity, error) {
	env := models.Environment{}
	err := repo.db.Unscoped().Where(criterias).First(&env).Error
	if gorm.IsRecordNotFoundError(err) {
		return &env, repositories.NewEntityDoesNotExistError(env, criterias)
	}
	return &env, err
}

// FindOneBy fetch the first application matching each criteria
func (repo *EnvironmentRepository) FindOneBy(criterias map[string]interface{}) (repositories.Entity, error) {
	env := models.Environment{}
	err := repo.db.Where(criterias).First(&env).Error
	if gorm.IsRecordNotFoundError(err) {
		return &env, repositories.NewEntityDoesNotExistError(env, criterias)
	}
	return &env, err
}

// Save persists an application to the database
func (repo *EnvironmentRepository) Save(environment repositories.Entity) error {
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
func (repo *EnvironmentRepository) Remove(env interface{}) error {
	env, err := repo.mustBeEntity(env)
	if err != nil {
		return err
	}

	return repo.db.Delete(env).Error
}

// FindAllPage returns a page of matching entities
func (repo *EnvironmentRepository) FindAllPage(pageable repositories.Pageable) (repositories.Page, error) {
	return repo.FindPageBy(pageable, map[string]interface{}{})
}

// FindPageBy returns a page of matching entities
func (repo *EnvironmentRepository) FindPageBy(pageable repositories.Pageable, criterias map[string]interface{}) (repositories.Page, error) {
	if pageable.Size == 0 {
		pageable.Size = environmentsDefaultPageSize
	}
	page := repositories.Page{Pageable: pageable, BasePath: models.EnvironmentBasePath}
	var environments []*models.Environment

	if err := repo.db.Where(criterias).Offset(pageable.Page * pageable.Size).Limit(pageable.Size).Find(&environments).Error; err != nil {
		return page, err
	}
	page.Content = environments

	count := 0
	if err := repo.db.Model(&models.Environment{}).Where(criterias).Count(&count).Error; err != nil {
		return page, err
	}
	page.TotalElements = count

	return page, nil
}

// Truncate empties the applications table for testing purposes
func (repo *EnvironmentRepository) Truncate() error {
	return repo.db.Delete(models.Environment{}).Error
}

func (repo *EnvironmentRepository) mustBeEntity(id interface{}) (*models.Environment, error) {
	var env *models.Environment
	switch id := id.(type) {
	case uint:
		env = &models.Environment{ID: id}
	case *models.Environment:
		env = id
	case models.Environment:
		env = &id
	default:
		return nil, repositories.NewUnsupportedEntityError(env, id)
	}
	return env, nil
}
