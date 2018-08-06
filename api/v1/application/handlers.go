package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/evanphx/json-patch"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/handlers"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
	"github.com/ovh/lhasa/api/v1/badge"
)

type badgeRatingsRequest struct {
	ApplicationDomain  string `path:"domain"`
	ApplicationName    string `path:"name"`
	ApplicationVersion string `path:"version"`
	BadgeSlug          string `path:"badgeslug"`
	Level              string `json:"level"`
	Comment            string `json:"comment"`
}

// HandlerCreate replace or create a resource
func HandlerCreate(repository *Repository, updater LatestUpdater) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, application *v1.Release) (*v1.Release, error) {
		oldres, err := repository.FindOneByUnscoped(map[string]interface{}{"domain": application.Domain, "name": application.Name, "version": application.Version})
		oldapp := oldres.(*v1.Release)
		if hateoas.IsEntityDoesNotExistError(err) {
			if err := updater(application, handlers.GetLogger(c)); err != nil {
				return nil, err
			}
			return nil, hateoas.ErrorCreated
		}
		if err != nil {
			return nil, err
		}

		application.ID = oldapp.ID
		application.CreatedAt = oldapp.CreatedAt

		if err := updater(application, handlers.GetLogger(c)); err != nil {
			return nil, err
		}
		if oldapp.DeletedAt != nil {
			return nil, hateoas.ErrorCreated
		}

		return application, nil

	}, http.StatusOK)
}

// HandlerPatch patch release properties
func HandlerPatch(repository *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *v1.Release) error {
		contentType := "application/merge-patch+json"
		if c.ContentType() != contentType {
			return errors.NotSupportedf("only application/merge-patch+json is supported, content-type %s", c.ContentType(), contentType)
		}
		entity, err := repository.FindOneBy(map[string]interface{}{"domain": request.Domain, "name": request.Name, "version": request.Version})
		if err != nil {
			return err
		}
		release, ok := entity.(*v1.Release)
		if !ok {
			return errors.New("invalid type")
		}

		if request.Properties == nil {
			return nil
		}

		document, err := release.Properties.MarshalJSON()
		if err != nil {
			return err
		}

		patch, err := json.Marshal(request.Properties)
		if err != nil {
			return err
		}

		properties, err := jsonpatch.MergePatch(document, patch)
		if err != nil {
			return err
		}

		if err := release.Properties.UnmarshalJSON(properties); err != nil {
			return err
		}

		return repository.Save(release)
	}, http.StatusNoContent)
}

func retrieveBadgeRatings(appv *v1.Release) (map[string]v1.BadgeRating, error) {
	badgeRatings := map[string]v1.BadgeRating{}
	if appv.BadgeRatings != nil {
		if len(appv.BadgeRatings.RawMessage) > 0 {
			if err := json.Unmarshal(appv.BadgeRatings.RawMessage, &badgeRatings); err != nil {
				return nil, err
			}
		}
	}
	return badgeRatings, nil
}

// HandlerGetBadgeRatingsForAppVersion retrieves all the badge values for a particular application version
func HandlerGetBadgeRatingsForAppVersion(repository *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *badgeRatingsRequest) ([]*v1.BadgeRating, error) {
		// Retrieve the release data
		appv, err := repository.FindOneByDomainNameVersion(request.ApplicationDomain, request.ApplicationName, request.ApplicationVersion)
		if err != nil {
			return nil, err
		}

		// Retrieve all the badges
		result := make([]*v1.BadgeRating, 0)
		badgeRepo := badge.NewRepository(repository.db)
		badges, err := badgeRepo.FindBy(map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		badgeRatings, err := retrieveBadgeRatings(appv)
		if err != nil {
			return nil, err
		}

		for _, bdg := range badges.([]*v1.Badge) {
			levels, err := badge.UnmarshalLevels(*bdg)
			if err != nil {
				return nil, err
			}
			rating := &v1.BadgeRating{
				Badge:      bdg,
				BadgeTitle: bdg.Title,
				Release:    appv,
				Value:      "unset",
				Comment:    "",
				Level:      &levels[0],
			}
			br, hasRating := badgeRatings[bdg.Slug]
			if hasRating {
				if level, found := searchLevel(levels, br); found {
					rating = &v1.BadgeRating{
						Badge:      bdg,
						Release:    appv,
						BadgeTitle: bdg.Title,
						Value:      br.Value,
						Comment:    br.Comment,
						Level:      &level,
					}
				}
			}
			rating.ToResource(hateoas.BaseURL(c))
			result = append(result, rating)
		}
		return result, nil
	}, http.StatusOK)
}

func searchLevel(levels []v1.BadgeLevel, br v1.BadgeRating) (v1.BadgeLevel, bool) {
	for _, badgeLevel := range levels {
		if badgeLevel.ID == br.Value {
			return badgeLevel, true
		}
	}
	return v1.BadgeLevel{}, false
}

// HandlerSetBadgeRatingForRelease sets the badge values for a particular release
func HandlerSetBadgeRatingForRelease(appRepo *Repository, badgeRepo *badge.Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *badgeRatingsRequest) (*v1.BadgeRating, error) {
		// Retrieve the application version data
		appv, err := appRepo.FindOneByDomainNameVersion(request.ApplicationDomain, request.ApplicationName, request.ApplicationVersion)

		if err != nil {
			return nil, err
		}
		badgeRatings, err := retrieveBadgeRatings(appv)
		if err != nil {
			return nil, err
		}

		// Retrieve the badge data
		bdg, err := badgeRepo.FindOneBySlug(request.BadgeSlug)
		if err != nil {
			return nil, err
		}

		_, err = badge.GetBadgeLevelByID(bdg, request.Level, false)
		if err != nil {
			return nil, err
		}

		badgeRatings[request.BadgeSlug] = v1.BadgeRating{
			BadgeID: request.BadgeSlug,
			Value:   request.Level,
			Comment: request.Comment,
		}
		badgeRatingsJSON, err := json.Marshal(badgeRatings)
		if err != nil {
			return nil, err
		}
		appv.BadgeRatings = &postgres.Jsonb{RawMessage: badgeRatingsJSON}
		br := &v1.BadgeRating{
			BadgeID: request.BadgeSlug,
			Badge:   bdg,
			Release: appv,
			Value:   request.Level,
			Comment: request.Comment,
		}
		br.ToResource(hateoas.BaseURL(c))
		return br, appRepo.Save(appv)
	}, http.StatusCreated)
}

// HandlerDeleteBadgeRatingForRelease sets the badge values for a particular release
func HandlerDeleteBadgeRatingForRelease(repository *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *badgeRatingsRequest) error {
		// Retrieve the application version data
		appv, err := repository.FindOneByDomainNameVersion(request.ApplicationDomain, request.ApplicationName, request.ApplicationVersion)
		if err != nil {
			return err
		}
		badgeRatings, err := retrieveBadgeRatings(appv)
		if err != nil {
			return err
		}

		delete(badgeRatings, request.BadgeSlug)

		badgeRatingsJSON, err := json.Marshal(badgeRatings)
		if err != nil {
			return err
		}
		appv.BadgeRatings = &postgres.Jsonb{RawMessage: badgeRatingsJSON}

		return repository.Save(appv)
	}, http.StatusOK)
}

// HandlerRedirectLatest reply back the latest version of this application
func HandlerRedirectLatest(appLatestRepo *RepositoryLatest) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, _ *v1.Application) (*v1.Release, error) {
		app, err := hateoas.FindByPath(c, appLatestRepo)
		if err != nil {
			return nil, err
		}
		if release, ok := app.(*v1.Release); ok {
			return release, nil
		}
		return nil, fmt.Errorf("returned entity %T is not a deployment", app)
	}, http.StatusOK)
}
