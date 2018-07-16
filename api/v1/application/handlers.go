package application

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/loopfz/gadgeto/tonic"
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
			if err := updater(application); err != nil {
				return nil, err
			}
			return nil, hateoas.ErrorCreated
		}
		if err != nil {
			return nil, err
		}

		application.ID = oldapp.ID
		application.CreatedAt = oldapp.CreatedAt

		if err := updater(application); err != nil {
			return nil, err
		}
		if oldapp.DeletedAt != nil {
			return nil, hateoas.ErrorCreated
		}

		return application, nil

	}, http.StatusOK)
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
		var result []*v1.BadgeRating
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
			unsetRating := &v1.BadgeRating{
				Badge:      bdg,
				BadgeTitle: bdg.Title,
				Release:    appv,
				Value:      "unset",
				Comment:    "",
				Level:      &levels[0],
			}
			unsetRating.ToResource(hateoas.BaseURL(c))
			b, found := badgeRatings[bdg.Slug]
			if !found {
				result = append(result, unsetRating)
				continue
			}

			if level, found := searchLevel(levels, b); found {
				rating := &v1.BadgeRating{
					Badge:      bdg,
					Release:    appv,
					BadgeTitle: bdg.Title,
					Value:      b.Value,
					Comment:    b.Comment,
					Level:      &level,
				}
				rating.ToResource(hateoas.BaseURL(c))
				result = append(result, rating)
				continue
			}

			result = append(result, unsetRating)
		}
		return result, nil
	}, http.StatusOK)
}

func searchLevel(levels []v1.BadgeLevel, b v1.BadgeRating) (v1.BadgeLevel, bool) {
	for _, badgeLevel := range levels {
		if badgeLevel.ID == b.Value {
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
		return &v1.BadgeRating{
			BadgeID: request.BadgeSlug,
			Value:   request.Level,
			Comment: request.Comment,
		}, appRepo.Save(appv)
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
			release.ToResource(hateoas.BaseURL(c))
			return release, nil
		}
		return nil, errors.New("internal type error")
	}, http.StatusOK)
}
