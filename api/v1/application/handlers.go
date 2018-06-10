package application

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
	"github.com/ovh/lhasa/api/v1/badge"
)

type badgeRating struct {
	BadgeSlug  string        `json:"badgeslug" validate:"not null; not empty"`
	BadgeTitle string        `json:"badgetitle"`
	Value      string        `json:"value" validate:"not null; not empty"`
	Comment    string        `json:"comment"`
	Level      v1.BadgeLevel `json:"level"`
}

type badgeRatingsRequest struct {
	ApplicationDomain  string `path:"domain"`
	ApplicationName    string `path:"name"`
	ApplicationVersion string `path:"version"`
	BadgeSlug          string `path:"badgeslug"`
	Level              string `json:"level"`
	Comment            string `json:"comment"`
}

// HandlerCreate replace or create a resource
func HandlerCreate(repository *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, application *v1.ApplicationVersion) error {
		oldres, err := repository.FindOneByUnscoped(map[string]interface{}{"domain": application.Domain, "name": application.Name, "version": application.Version})
		oldapp := oldres.(*v1.ApplicationVersion)
		if hateoas.IsEntityDoesNotExistError(err) {
			if err := repository.Save(application); err != nil {
				return err
			}
			return hateoas.ErrorCreated
		}
		if err != nil {
			return err
		}

		application.ID = oldapp.ID
		application.CreatedAt = oldapp.CreatedAt
		if err := repository.Save(application); err != nil {
			return err
		}
		if oldapp.DeletedAt != nil {
			return hateoas.ErrorCreated
		}
		return nil

	}, http.StatusOK)
}

func retrieveBadgeRatings(appv *v1.ApplicationVersion) ([]badgeRating, error) {
	badgeRatings := make([]badgeRating, 0)
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
	return tonic.Handler(func(c *gin.Context, request *badgeRatingsRequest) (*[]badgeRating, error) {

		// Retrieve the application version data
		appv, err := repository.FindOneByDomainNameVersion(request.ApplicationDomain, request.ApplicationName, request.ApplicationVersion)
		if err != nil {
			return nil, err
		}

		// Retrieve all the badges
		result := make([]badgeRating, 0)
		badgeRepo := badge.NewRepository(repository.db)
		badges, err := badgeRepo.FindAll()
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

			found := false
			for _, bdgRating := range badgeRatings {
				if bdgRating.BadgeSlug != bdg.Slug {
					continue
				}
				for _, badgeLevel := range levels {
					if badgeLevel.ID == bdgRating.Value {
						result = append(result, badgeRating{
							BadgeSlug:  bdg.Slug,
							BadgeTitle: bdg.Title,
							Value:      bdgRating.Value,
							Comment:    bdgRating.Comment,
							Level:      badgeLevel,
						})
						found = true
						break
					}
				}

			}
			if found == false {
				result = append(result, badgeRating{
					BadgeSlug:  bdg.Slug,
					BadgeTitle: bdg.Title,
					Value:      "unset",
					Comment:    "",
					Level:      levels[0],
				})
			}
		}
		return &result, nil
	}, http.StatusOK)
}

// HandlerSetBadgeRatingForAppVersion sets the badge values for a particular application version
func HandlerSetBadgeRatingForAppVersion(repository *Repository) gin.HandlerFunc {
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

		// Retrieve the badge data
		badgeRepo := badge.NewRepository(repository.db)
		bdg, err := badgeRepo.FindOneBySlug(request.BadgeSlug)
		if err != nil {
			return err
		}

		_, err = badge.GetBadgeLevelByID(bdg, request.Level, false)
		if err != nil {
			return err
		}
		found := false
		for i, bdgRating := range badgeRatings {
			if bdgRating.BadgeSlug == request.BadgeSlug {
				badgeRatings[i] = badgeRating{
					BadgeSlug: request.BadgeSlug,
					Value:     request.Level,
					Comment:   request.Comment,
				}
				found = true
				break
			}
		}
		if found == false {
			badgeRatings = append(badgeRatings, badgeRating{
				BadgeSlug: request.BadgeSlug,
				Value:     request.Level,
				Comment:   request.Comment,
			})
		}
		badgeRatingsJSON, err := json.Marshal(badgeRatings)
		if err != nil {
			return err
		}
		appv.BadgeRatings = &postgres.Jsonb{RawMessage: badgeRatingsJSON}
		return repository.Save(appv)
	}, http.StatusCreated)
}

// HandlerDeleteBadgeRatingForAppVersion sets the badge values for a particular application version
func HandlerDeleteBadgeRatingForAppVersion(repository *Repository) gin.HandlerFunc {
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

		for i, bdgRating := range badgeRatings {
			if bdgRating.BadgeSlug == request.BadgeSlug {
				badgeRatings = append(badgeRatings[:i], badgeRatings[i+1:]...)
				break
			}
		}

		badgeRatingsJSON, err := json.Marshal(badgeRatings)
		if err != nil {
			return err
		}
		appv.BadgeRatings = &postgres.Jsonb{RawMessage: badgeRatingsJSON}

		return repository.Save(appv)
	}, http.StatusOK)
}
