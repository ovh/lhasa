package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ovh/lhasa/api/handlers"
	"github.com/ovh/lhasa/api/repositories"
	"github.com/ovh/lhasa/api/v1/models"
	apprepo "github.com/ovh/lhasa/api/v1/repositories"
)

// ApplicationCreate replace or create a resource
func ApplicationCreate(repository *apprepo.ApplicationRepository) func(*gin.Context, *models.Application) error {
	return func(c *gin.Context, application *models.Application) error {
		oldapp, err := repository.FindOneByUnscoped(map[string]interface{}{"domain": application.Domain, "name": application.Name, "version": application.Version})
		if !repositories.IsEntityDoesNotExistError(err) {
			application.ID = oldapp.ID
			application.CreatedAt = oldapp.CreatedAt
			if err := repository.Save(application); err != nil {
				return err
			}
			if oldapp.DeletedAt != nil {
				return handlers.RestErrorCreated
			}
			return nil
		}
		if err := repository.Save(application); err != nil {
			return err
		}
		return handlers.RestErrorCreated
	}
}
