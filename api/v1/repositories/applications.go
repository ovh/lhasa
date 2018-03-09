package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/ovh/lhasa/api/repositories"
	"github.com/ovh/lhasa/api/v1/models"
)

const (
	applicationsDefaultPageSize = 20
)

// ApplicationRepository is a repository manager for applications
type ApplicationRepository struct {
	db      *gorm.DB
	countDB *gorm.DB
}

// NewApplicationVersionAwareRepository creates a version aware application repository
func NewApplicationVersionAwareRepository(db *gorm.DB) *ApplicationRepository {
	repository := NewApplicationRepository(
		db.Group("id, name, domain").Having("max(id) = id"),
	)
	repository.countDB = db.Group("name, domain")
	return repository
}

// NewApplicationRepository creates an application repository
func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{
		db:      db,
		countDB: db,
	}
}

// Migrate run automated gorm migrations
func (repo *ApplicationRepository) Migrate() error {
	if err := repo.db.AutoMigrate(models.Application{}).Error; err != nil {
		return err
	}
	if err := repo.db.AutoMigrate(models.PersonInfo{}).Error; err != nil {
		return err
	}
	return repo.db.AutoMigrate(models.DependencyInfo{}).Error
}

// FindAll returns all entities of the repository type
func (repo *ApplicationRepository) FindAll() (interface{}, error) {
	return repo.FindBy(map[string]interface{}{})
}

// FindAllPage returns a page of matching entities
func (repo *ApplicationRepository) FindAllPage(pageable repositories.Pageable) (repositories.Page, error) {
	return repo.FindPageBy(pageable, map[string]interface{}{})
}

// FindPageBy returns a page of matching entities
func (repo *ApplicationRepository) FindPageBy(pageable repositories.Pageable, criterias map[string]interface{}) (repositories.Page, error) {
	if pageable.Size == 0 {
		pageable.Size = applicationsDefaultPageSize
	}
	page := repositories.Page{Pageable: pageable}
	var applications []models.Application

	if err := repo.db.Where(criterias).Offset(pageable.Page * pageable.Size).Limit(pageable.Size).Find(&applications).Error; err != nil {
		return page, err
	}
	page.Content = applications

	count := 0
	if err := repo.countDB.Model(&models.Application{}).Where(criterias).Count(&count).Error; err != nil {
		return page, err
	}
	page.TotalElements = count

	if pageable.IndexedBy != "" {
		currentIndex := map[string][]models.Application{}
		ids := map[string]bool{}
		for _, application := range applications {
			indexedField, err := getIndexedField(pageable.IndexedBy, application)
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

func getIndexedField(field string, application models.Application) (string, error) {
	switch field {
	case "version":
		return application.Version, nil
	case "profile":
	case "domain":
		return application.Domain, nil
	case "type":
		return application.Type, nil
	}
	return "", repositories.NewUnsupportedIndexError(field, "version", "domain", "type")
}

// Save persists an application to the database
func (repo *ApplicationRepository) Save(application interface{}) error {
	app, err := mustBeApplication(application)
	if err != nil {
		return err
	}

	oldapp, err := repo.findOneByUnscoped(map[string]interface{}{"domain": app.Domain, "name": app.Name, "version": app.Version})
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if gorm.IsRecordNotFoundError(err) {
		return repo.db.Create(app).Error
	}
	app.ID = oldapp.ID
	return repo.db.Unscoped().Save(app).Error
}

// Truncate empties the applications table for testing purposes
func (repo *ApplicationRepository) Truncate() error {
	return repo.db.Delete(models.Application{}).Error
}

// Remove deletes the application whose ID is given as a parameter
func (repo *ApplicationRepository) Remove(app interface{}) error {
	app, err := mustBeApplication(app)
	if err != nil {
		return err
	}

	return repo.db.Delete(app).Error
}

// FindByID gives the details of a particular application
func (repo *ApplicationRepository) FindByID(id interface{}) (interface{}, error) {
	app := models.Application{}
	if err := repo.db.First(&app, id).Error; err != nil {
		return nil, err
	}
	return app, nil
}

// findOneByUnscoped gives the details of a particular application, even if soft deleted
func (repo *ApplicationRepository) findOneByUnscoped(criterias map[string]interface{}) (models.Application, error) {
	app := models.Application{}
	err := repo.db.Unscoped().Where(criterias).First(&app).Error
	return app, err
}

// FindBy fetch a collection of applications matching each criteria
func (repo *ApplicationRepository) FindBy(criterias map[string]interface{}) (interface{}, error) {
	var apps []models.Application
	err := repo.db.Where(criterias).Find(&apps).Error
	return apps, err
}

// FindOneBy fetch the first application matching each criteria
func (repo *ApplicationRepository) FindOneBy(criterias map[string]interface{}) (interface{}, error) {
	app := models.Application{}
	err := repo.db.Where(criterias).First(&app).Error
	if gorm.IsRecordNotFoundError(err) {
		return app, repositories.NewEntityDoesNotExistError(app, criterias)
	}
	return app, err
}

func mustBeApplication(id interface{}) (*models.Application, error) {
	var app *models.Application
	switch id := id.(type) {
	case uint:
		app = &models.Application{ID: id}
	case *models.Application:
		app = id
	case models.Application:
		app = &id
	default:
		return nil, repositories.NewUnsupportedEntityError(app, id)
	}
	return app, nil
}
