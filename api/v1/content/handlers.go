package content

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
)

// HandlerGet returns the first resource matching path params
func HandlerGet(repository *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, content *v1.Content) (interface{}, error) {
		if len(content.Locale) == 0 {
			content.Locale = "en-GB"
		}
		result, err := repository.FindOneByUnscoped(map[string]interface{}{"name": content.Name, "locale": content.Locale})
		if hateoas.IsEntityDoesNotExistError(err) {
			if err := repository.Save(content); err != nil {
				return nil, err
			}
			return nil, hateoas.EntityDoesNotExistError{
				Criteria:   map[string]interface{}{"name": content.Name, "locale": content.Locale},
				EntityName: "Content",
			}
		}

		return result, nil
	}, http.StatusOK)
}

// HandlerCreate replace or create a resource
func HandlerCreate(repository *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, content *v1.Content) error {
		oldres, err := repository.FindOneByUnscoped(map[string]interface{}{"name": content.Name, "locale": content.Locale})
		oldapp := oldres.(*v1.Content)
		if hateoas.IsEntityDoesNotExistError(err) {
			if err := repository.Save(content); err != nil {
				return err
			}
			return hateoas.ErrorCreated
		}
		if err != nil {
			return err
		}

		content.ID = oldapp.ID
		content.CreatedAt = oldapp.CreatedAt
		if err := repository.Save(content); err != nil {
			return err
		}
		if oldapp.DeletedAt != nil {
			return hateoas.ErrorCreated
		}
		c.Data(http.StatusOK, content.ContentType, content.Body)
		return nil

	}, http.StatusOK)
}
