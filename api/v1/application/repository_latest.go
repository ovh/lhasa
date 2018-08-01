package application

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
)

const querySelectLatestApplications = `SELECT "av".* FROM "releases" as "av"
JOIN "applications" ON "applications"."latest_release_id" = "av"."id"
WHERE "av"."deleted_at" IS NULL
%s
ORDER BY %s
LIMIT %d OFFSET %d`

const querySelectLatestApplicationsCount = `SELECT count(1) as nbApp FROM "releases" as "av"
JOIN "applications" ON "applications"."latest_release_id" = "av"."id"
WHERE "av"."deleted_at" IS NULL
%s`

const querySelectLatestApplicationsUnscoped = `SELECT "av".* FROM "releases" as "av"
JOIN "applications" ON "applications"."latest_release_id" = "av"."id"
WHERE 1 = 1
%s
ORDER BY %s
LIMIT %d OFFSET %d`

// RepositoryLatest is a repository manager for applications
type RepositoryLatest struct {
	db *gorm.DB
	Repository
}

// RepositoryLatestFactory defines a repository constructor
type RepositoryLatestFactory func(*gorm.DB) *RepositoryLatest

// NewRepositoryLatest creates an application repository
func NewRepositoryLatest(db *gorm.DB) *RepositoryLatest {
	return &RepositoryLatest{
		db: db,
	}
}

// FindPageBy returns a page of matching entities
func (repo *RepositoryLatest) FindPageBy(pageable hateoas.Pageable, criteria map[string]interface{}) (hateoas.Page, error) {
	page := hateoas.NewPage(pageable, defaultPageSize, v1.ApplicationBasePath)
	var applications []*v1.Release

	whereClause := getWhereClause(criteria)
	if err := repo.db.
		Raw(fmt.Sprintf(querySelectLatestApplications,
			whereClause,
			page.Pageable.GetSortClause(),
			page.Pageable.Size,
			page.Pageable.GetOffset())).
		Scan(&applications).Error; err != nil {
		return page, err
	}
	page.Content = applications

	count := 0
	rows, err := repo.db.Raw(fmt.Sprintf(querySelectLatestApplicationsCount, whereClause)).Rows()
	if err != nil {
		return page, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&count)
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

func getWhereClause(criteria map[string]interface{}) string {
	whereClause := ""
	for key, value := range criteria {
		if value == "" {
			continue
		}
		whereClause += getWhereClauseChunk(key, value.(string))
	}
	return whereClause
}

func (repo *RepositoryLatest) getIndexedField(field string, application *v1.Release) (string, error) {
	switch field {
	case "profile":
	case "domain":
		return application.Domain, nil
	}
	return "", hateoas.NewUnsupportedIndexError(field, "domain")
}

func getWhereClauseChunk(field string, value string) string {
	// TODO FIX SQL injection issue
	if field == v1.FreeSearch {
		return fmt.Sprintf(
			`AND ( "av".name ILIKE '%%%s%%'
				OR "av".domain ILIKE '%%%s%%'
				OR "av".version ILIKE '%%%s%%'
				OR cast("av".properties as TEXT) ILIKE '%%%s%%'
				OR cast("av".manifest as TEXT) ILIKE '%%%s%%'
				OR cast("av".tags as TEXT) ILIKE '%%%s%%'
			) `, value, value, value, value, value, value)
	}
	if !strings.Contains(field, ".") {
		return fmt.Sprintf("AND \"av\".%q = '%s' ", field, value)
	}
	fieldParts := strings.Split(field, ".")
	clause := ""
	for i, p := range fieldParts {
		switch i {
		case 0:
			clause = fmt.Sprintf("AND \"av\".%q::json", p)
		case len(fieldParts) - 1:
			clause += fmt.Sprintf("->>'%s'", p)
		default:
			clause += fmt.Sprintf("->'%s'", p)
		}
	}
	return clause + fmt.Sprintf("= '%s'", value)
}

// Save persists an application to the database
func (repo *RepositoryLatest) Save(application hateoas.Entity) error {
	app, err := repo.mustBeEntity(application)
	if err != nil {
		return err
	}

	if app.ID == 0 {
		return repo.db.Create(app).Error
	}
	return repo.db.Unscoped().Save(app).Error
}

// FindOneByUnscoped gives the details of a particular application, even if soft deleted
func (repo *RepositoryLatest) FindOneByUnscoped(criteria map[string]interface{}) (hateoas.SoftDeletableEntity, error) {
	app := v1.Release{}
	err := repo.db.Raw(fmt.Sprintf(querySelectLatestApplicationsUnscoped,
		getWhereClause(criteria), "1", 1, 0)).Scan(&app).Error
	if gorm.IsRecordNotFoundError(err) {
		return &app, hateoas.NewEntityDoesNotExistError(app, criteria)
	}
	return &app, err
}

// FindBy fetch a collection of applications matching each criteria
func (repo *RepositoryLatest) FindBy(criteria map[string]interface{}) (interface{}, error) {
	var apps []*v1.Application
	err := repo.db.Where(criteria).Find(&apps).Error
	return apps, err
}

// FindOneBy find by criteria
func (repo *RepositoryLatest) FindOneBy(criteria map[string]interface{}) (hateoas.Entity, error) {
	app := v1.Release{}
	err := repo.db.Raw(
		fmt.Sprintf(querySelectLatestApplications, getWhereClause(criteria), "1", 1, 0)).
		Scan(&app).Error
	if gorm.IsRecordNotFoundError(err) {
		return &app, hateoas.NewEntityDoesNotExistError(app, criteria)
	}
	return &app, err
}

// FindApplication finds a single application
func (repo *RepositoryLatest) FindApplication(domain, name string) (*v1.Application, error) {
	app := v1.Application{}
	criteria := map[string]interface{}{"domain": domain, "name": name}
	err := repo.db.Where(criteria).Preload("LatestRelease").First(&app).Error
	if gorm.IsRecordNotFoundError(err) {
		return &app, hateoas.NewEntityDoesNotExistError(app, criteria)
	}
	return &app, err
}

// FindOneByDomainNameVersion fetch the first application matching each criteria
func (repo *RepositoryLatest) FindOneByDomainNameVersion(domain, name, version string) (*v1.Release, error) {
	app := v1.Release{}
	criteria := map[string]interface{}{
		"domain": domain,
		"name":   name,
	}
	err := repo.db.Raw(
		// Order by first column, limit 1, offset 0
		fmt.Sprintf(querySelectLatestApplications, getWhereClause(criteria), "1", 1, 0)).
		Scan(&app).Error
	if gorm.IsRecordNotFoundError(err) {
		return &app, hateoas.NewEntityDoesNotExistError(app, criteria)
	}
	return &app, err
}

func (repo *RepositoryLatest) mustBeEntity(id interface{}) (*v1.Application, error) {
	var app *v1.Application
	switch id := id.(type) {
	case uint:
		app = &v1.Application{ID: id}
	case *v1.Application:
		app = id
	case v1.Application:
		app = &id
	default:
		return nil, hateoas.NewUnsupportedEntityError(app, id)
	}
	return app, nil
}
