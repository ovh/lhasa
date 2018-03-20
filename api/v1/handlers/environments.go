package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/handlers"
	"github.com/ovh/lhasa/api/repositories"
	"github.com/ovh/lhasa/api/v1/models"
	envrepo "github.com/ovh/lhasa/api/v1/repositories"
)

// EnvironmentCreate replace or create a resource
func EnvironmentCreate(repository *envrepo.EnvironmentRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, env *models.Environment) error {
		oldres, err := repository.FindOneByUnscoped(map[string]interface{}{"slug": env.Slug})
		oldenv := oldres.(*models.Environment)
		if repositories.IsEntityDoesNotExistError(err) {
			if err := repository.Save(env); err != nil {
				return err
			}
			return handlers.RestErrorCreated
		}
		if err != nil {
			return err
		}
		env.ID = oldenv.ID
		env.CreatedAt = oldenv.CreatedAt
		if err := repository.Save(env); err != nil {
			return err
		}
		if oldenv.DeletedAt != nil {
			return handlers.RestErrorCreated
		}
		return nil
	}, http.StatusOK)
}
