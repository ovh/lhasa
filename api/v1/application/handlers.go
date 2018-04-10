package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
)

// HandlerCreate replace or create a resource
func HandlerCreate(repository *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, application *v1.Application) error {
		oldres, err := repository.FindOneByUnscoped(map[string]interface{}{"domain": application.Domain, "name": application.Name, "version": application.Version})
		oldapp := oldres.(*v1.Application)
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
