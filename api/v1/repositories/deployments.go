package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/ovh/lhasa/api/repositories"
	"github.com/ovh/lhasa/api/v1/models"
)

const (
	deploymentsDefaultPageSize = 20

	queryUndeployByApplicationNameEnvSlug = "UPDATE \"deployments\" AS \"d\" " +
		"SET \"undeployed_at\" = now(), \"updated_at\" = now() " +
		"FROM \"applications\" as \"a\" " +
		"WHERE \"d\".\"deleted_at\" IS NULL " +
		"AND \"a\".\"id\" = \"d\".\"application_id\" " +
		"AND \"a\".\"domain\" = ? " +
		"AND \"a\".\"name\" = ? " +
		"AND \"d\".\"environment_id\" = ? " +
		"AND \"undeployed_at\" IS NULL"
)

// DeploymentRepository is a repository manager for applications
type DeploymentRepository struct {
	db *gorm.DB
}

// NewDeploymentRepository creates an application repository
func NewDeploymentRepository(db *gorm.DB) *DeploymentRepository {
	return &DeploymentRepository{
		db: db,
	}
}

// GetNewEntityInstance returns a new empty instance of the entity managed by this repository
func (repo *DeploymentRepository) GetNewEntityInstance() repositories.Entity {
	return &models.Deployment{}
}

// FindAllPage returns a page of matching entities
func (repo *DeploymentRepository) FindAllPage(pageable repositories.Pageable) (repositories.Page, error) {
	return repo.FindPageBy(pageable, map[string]interface{}{})
}

// FindPageBy returns a page of matching entities
func (repo *DeploymentRepository) FindPageBy(pageable repositories.Pageable, criterias map[string]interface{}) (repositories.Page, error) {
	if pageable.Size == 0 {
		pageable.Size = deploymentsDefaultPageSize
	}
	page := repositories.Page{Pageable: pageable, BasePath: models.DeploymentBasePath}
	var deployments []*models.Deployment

	if err := repo.db.Preload("Environment").Preload("Application").Model(models.Deployment{}).
		Offset(pageable.Page*pageable.Size).Limit(pageable.Size).Find(&deployments, criterias).Error; err != nil {
		return page, err
	}
	page.Content = deployments

	count := 0
	if err := repo.db.Model(models.Deployment{}).Where(criterias).Count(&count).Error; err != nil {
		return page, err
	}
	page.TotalElements = count

	return page, nil
}

// FindAll returns all entities of the repository type
func (repo *DeploymentRepository) FindAll() (interface{}, error) {
	return repo.FindBy(map[string]interface{}{})
}

// Save persists an deployment to the database
func (repo *DeploymentRepository) Save(deployment repositories.Entity) error {
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
func (repo *DeploymentRepository) Truncate() error {
	return repo.db.Delete(models.Deployment{}).Error
}

// Remove deletes the deployment whose GetID is given as a parameter
func (repo *DeploymentRepository) Remove(dep interface{}) error {
	dep, err := repo.mustBeEntity(dep)
	if err != nil {
		return err
	}

	return repo.db.Delete(dep).Error
}

// FindByID gives the details of a particular deployment
func (repo *DeploymentRepository) FindByID(id interface{}) (repositories.Entity, error) {
	dep := models.Deployment{}
	if err := repo.db.First(&dep, id).Error; err != nil {
		return nil, err
	}
	return &dep, nil
}

// FindOneByUnscoped gives the details of a particular deployment, even if soft deleted
func (repo *DeploymentRepository) FindOneByUnscoped(criterias map[string]interface{}) (repositories.SoftDeletableEntity, error) {
	dep := models.Deployment{}
	err := repo.db.Model(models.Deployment{}).Unscoped().First(dep, criterias).Error
	if gorm.IsRecordNotFoundError(err) {
		return &dep, repositories.NewEntityDoesNotExistError(dep, criterias)
	}
	return &dep, err
}

// FindBy fetch a collection of deployments matching each criteria
func (repo *DeploymentRepository) FindBy(criterias map[string]interface{}) (interface{}, error) {
	var deps []*models.Deployment
	err := repo.db.Model(models.Deployment{}).Find(&deps, criterias).Error
	return deps, err
}

// FindActivesBy fetch a collection of deployments matching each criteria
func (repo *DeploymentRepository) FindActivesBy(criterias map[string]interface{}) ([]*models.Deployment, error) {
	var deps []*models.Deployment

	err := repo.db.Preload("Environment").Preload("Application").Model(models.Deployment{}).Where("undeployed_at IS NULL").Find(&deps, criterias).Error
	return deps, err
}

// FindOneBy fetch the first deployment matching each criteria
func (repo *DeploymentRepository) FindOneBy(criterias map[string]interface{}) (repositories.Entity, error) {
	dep := models.Deployment{}
	err := repo.db.Where(criterias).First(&dep).Error
	if gorm.IsRecordNotFoundError(err) {
		return &dep, repositories.NewEntityDoesNotExistError(dep, criterias)
	}
	return &dep, err
}

// UndeployByApplicationEnv updates all deployments attached to a given application regardless version with an undeploy date to now
func (repo *DeploymentRepository) UndeployByApplicationEnv(domain, name string, envID uint) error {
	return repo.db.Exec(queryUndeployByApplicationNameEnvSlug, domain, name, envID).Error
}

func (repo *DeploymentRepository) mustBeEntity(id interface{}) (*models.Deployment, error) {
	var dep *models.Deployment
	switch id := id.(type) {
	case uint:
		dep = &models.Deployment{ID: id}
	case *models.Deployment:
		dep = id
	case models.Deployment:
		dep = &id
	default:
		return nil, repositories.NewUnsupportedEntityError(dep, id)
	}
	return dep, nil
}
