package environment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
)

// HandlerCreate replace or create a resource
func HandlerCreate(repository *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, env *v1.Environment) error {
		oldres, err := repository.FindOneByUnscoped(map[string]interface{}{"slug": env.Slug})
		oldenv := oldres.(*v1.Environment)
		if hateoas.IsEntityDoesNotExistError(err) {
			if err := repository.Save(env); err != nil {
				return err
			}
			return hateoas.ErrorCreated
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
			return hateoas.ErrorCreated
		}
		return nil
	}, http.StatusOK)
}
