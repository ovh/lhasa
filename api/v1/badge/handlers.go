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
	return tonic.Handler(func(c *gin.Context, bdg *v1.Badge) (*v1.Badge, error) {
		_, err := GetDefaultLevel(bdg)
		if err != nil {
			return nil, err
		}
		oldres, err := repository.FindOneByUnscoped(map[string]interface{}{"slug": bdg.Slug})
		if hateoas.IsEntityDoesNotExistError(err) {
			if err := repository.Save(bdg); err != nil {
				return nil, err
			}
			return nil, hateoas.ErrorCreated
		}
		if err != nil {
			return nil, err
		}
		oldbdg := oldres.(*v1.Badge)
		bdg.ID = oldbdg.ID
		bdg.CreatedAt = oldbdg.CreatedAt
		if err := repository.Save(bdg); err != nil {
			return nil, err
		}
		if oldbdg.DeletedAt != nil {
			return bdg, hateoas.ErrorCreated
		}
		return bdg, nil
	}, http.StatusOK)
}

// HandlerStats computes and renders badge statistics
func HandlerStats(repo *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, b *v1.Badge) (map[string]int, error) {
		bdg, err := repo.FindOneBySlug(b.Slug)
		if err != nil {
			return nil, err
		}
		l, err := GetDefaultLevel(bdg)
		if err != nil {
			return nil, err
		}
		return repo.GatherStats(b.Slug, l.ID)
	}, http.StatusOK)
}
