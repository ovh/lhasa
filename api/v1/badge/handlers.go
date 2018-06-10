package badge

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
)

// HandlerCreate replace or create a resource
func HandlerCreate(repository *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, bdg *v1.Badge) error {
		_, err := GetDefaultLevel(bdg)
		if err != nil {
			return err
		}
		oldres, err := repository.FindOneByUnscoped(map[string]interface{}{"slug": bdg.Slug})
		if hateoas.IsEntityDoesNotExistError(err) {
			if err := repository.Save(bdg); err != nil {
				return err
			}
			return hateoas.ErrorCreated
		}
		if err != nil {
			return err
		}
		oldbdg := oldres.(*v1.Badge)
		bdg.ID = oldbdg.ID
		bdg.CreatedAt = oldbdg.CreatedAt
		if err := repository.Save(bdg); err != nil {
			return err
		}
		if oldbdg.DeletedAt != nil {
			return hateoas.ErrorCreated
		}
		return nil
	}, http.StatusOK)
}
