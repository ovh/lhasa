package badge

import (
	"encoding/json"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
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
	return reflect.TypeOf(v1.Badge{})
}

// GetNewEntityInstance returns a new empty instance of the entity managed by this repository
func (repo *Repository) GetNewEntityInstance() hateoas.Entity {
	return &v1.Badge{}
}

// FindAll returns all entities of the repository type
func (repo *Repository) FindAll() (interface{}, error) {
	return repo.FindBy(map[string]interface{}{})
}

// FindByID gives the details of a particular application
func (repo *Repository) FindByID(id interface{}) (hateoas.Entity, error) {
	badge := v1.Badge{}
	err := repo.db.First(&badge, id).Error
	if gorm.IsRecordNotFoundError(err) {
		return &badge, hateoas.NewEntityDoesNotExistError(badge, map[string]interface{}{"id": id})
	}
	if err != nil {
		return nil, err
	}
	return &badge, nil
}

// FindOneBySlug fetch a collection of applications matching each criteria
func (repo *Repository) FindOneBySlug(slug string) (*v1.Badge, error) {
	badge := v1.Badge{}
	criterias := map[string]interface{}{"slug": slug}
	err := repo.db.First(&badge, criterias).Error
	if gorm.IsRecordNotFoundError(err) {
		return &badge, hateoas.NewEntityDoesNotExistError(badge, criterias)
	}
	if err != nil {
		return nil, err
	}
	return &badge, nil
}

// FindBy fetch a collection of applications matching each criteria
func (repo *Repository) FindBy(criterias map[string]interface{}) (interface{}, error) {
	var badges []*v1.Badge
	err := repo.db.Where(criterias).Find(&badges).Error
	return badges, err
}

// FindOneByUnscoped gives the details of a particular badge, even if soft deleted
func (repo *Repository) FindOneByUnscoped(criterias map[string]interface{}) (hateoas.SoftDeletableEntity, error) {
	badge := v1.Badge{}
	err := repo.db.Unscoped().Where(criterias).First(&badge).Error
	if gorm.IsRecordNotFoundError(err) {
		return &badge, hateoas.NewEntityDoesNotExistError(badge, criterias)
	}
	return &badge, err
}

// FindOneBy fetch the first badge matching each criteria
func (repo *Repository) FindOneBy(criterias map[string]interface{}) (hateoas.Entity, error) {
	badge := v1.Badge{}
	err := repo.db.Where(criterias).First(&badge).Error
	if gorm.IsRecordNotFoundError(err) {
		return &badge, hateoas.NewEntityDoesNotExistError(badge, criterias)
	}
	return &badge, err
}

// UnmarshalLevels extracts the levels JSONB field form the badge
func UnmarshalLevels(badge v1.Badge) ([]v1.BadgeLevel, error) {
	var levels []v1.BadgeLevel
	if err := json.Unmarshal(badge.Levels.RawMessage, &levels); err != nil {
		return nil, errors.BadRequestf("levels field should be a valid json object: %s", err.Error())
	}
	return levels, nil

}

// GetDefaultLevel returns the default level of a badge
func GetDefaultLevel(bdg *v1.Badge) (*v1.BadgeLevel, error) {
	levels, err := UnmarshalLevels(*bdg)
	if err != nil {
		return nil, err
	}
	var found *v1.BadgeLevel
	for _, lvl := range levels {
		if lvl.IsDefault {
			if found != nil {
				// We have two default levels -> error
				return nil, errors.BadRequestf("there are more than one levels with `isdefault=true`")
			}
			found = &lvl
		}
	}
	if found == nil {
		// We have two default levels -> error
		return nil, errors.BadRequestf("there is no level with `isdefault=true`")
	}
	return found, nil
}

// GetBadgeLevelByID returns the level identified by id
func GetBadgeLevelByID(bdg *v1.Badge, id string, returnDefaultIfNotExists bool) (*v1.BadgeLevel, error) {
	levels, err := UnmarshalLevels(*bdg)
	if err != nil {
		return nil, err
	}
	var defaultLevel *v1.BadgeLevel
	for _, lvl := range levels {
		if lvl.ID == id {
			return &lvl, nil
		}
		if lvl.IsDefault {
			defaultLevel = &lvl
		}
	}
	if !returnDefaultIfNotExists {
		return nil, errors.NotFoundf("there is no level id `%s`", id)
	}
	if defaultLevel == nil {
		return nil, errors.Errorf("there is no level with `isdefault=true`")
	}
	return defaultLevel, nil
}

// Save persists a badge to the database
func (repo *Repository) Save(badge hateoas.Entity) error {
	bdg, err := repo.mustBeEntity(badge)
	if err != nil {
		return err
	}
	levels, err := UnmarshalLevels(*bdg)
	if err != nil {
		return err
	}
	if levels[0].ID != "unset" {
		return errors.New("the `unset` level is mandatory")
	}
	if bdg.ID == 0 {
		return repo.db.Create(badge).Error
	}
	return repo.db.Unscoped().Save(badge).Error
}

// Remove deletes the application whose GetID is given as a parameter
func (repo *Repository) Remove(badge interface{}) error {
	bdg, err := repo.mustBeEntity(badge)
	if err != nil {
		return err
	}

	return repo.db.Delete(bdg).Error
}

// FindAllPage returns a page of matching entities
func (repo *Repository) FindAllPage(pageable hateoas.Pageable) (hateoas.Page, error) {
	return repo.FindPageBy(pageable, map[string]interface{}{})
}

// FindPageBy returns a page of matching entities
func (repo *Repository) FindPageBy(pageable hateoas.Pageable, criterias map[string]interface{}) (hateoas.Page, error) {
	page := hateoas.NewPage(pageable, defaultPageSize, v1.BadgeBasePath)
	var badges []*v1.Badge

	if err := repo.db.Where(criterias).
		Order(page.Pageable.GetSortClause()).
		Limit(page.Pageable.Size).
		Offset(page.Pageable.GetOffset()).
		Find(&badges).Error; err != nil {
		return page, err
	}
	page.Content = badges

	count := 0
	if err := repo.db.Model(&v1.Badge{}).Where(criterias).Count(&count).Error; err != nil {
		return page, err
	}
	page.TotalElements = count

	return page, nil
}

// Truncate empties the applications table for testing purposes
func (repo *Repository) Truncate() error {
	return repo.db.Delete(v1.Badge{}).Error
}

func (repo *Repository) mustBeEntity(id interface{}) (*v1.Badge, error) {
	var badge *v1.Badge
	switch id := id.(type) {
	case uint:
		badge = &v1.Badge{ID: id}
	case *v1.Badge:
		badge = id
	case v1.Badge:
		badge = &id
	default:
		return nil, hateoas.NewUnsupportedEntityError(badge, id)
	}
	return badge, nil
}
