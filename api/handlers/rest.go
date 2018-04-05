package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/repositories"
)

type restErrorCreated string
type restErrorGone string

// Error implements error interface
func (err restErrorCreated) Error() string {
	return string(err)
}

// Error implements error interface
func (err restErrorGone) Error() string {
	return string(err)
}

// RestErrorCreated is raised when no error occurs but a resource has been created (tonic single-status code workaround)
var RestErrorCreated = restErrorCreated("created")

// RestErrorGone is raised when a former resource has been requested but no longer exist
var RestErrorGone = restErrorGone("gone")

// RestFindByPage returns a filtered and paginated resource list
func RestFindByPage(repository repositories.PageableRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (*repositories.PagedResources, error) {
		pageable := repositories.Pageable{}
		c.ShouldBindQuery(&pageable)

		results, err := repository.FindPageBy(pageable, parsePathParams(c))
		if err != nil {
			return nil, err
		}

		resources := results.ToResources(HateoasBaseURL(c))
		return &resources, nil
	}, http.StatusPartialContent)
}

// RestFindByIndexedPage returns a filtered and paginated resource list, forcing an indexation
func RestFindByIndexedPage(repository repositories.PageableRepository, indexField string) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (*repositories.PagedResources, error) {
		pageable := repositories.Pageable{}
		c.ShouldBindQuery(&pageable)
		pageable.IndexedBy = indexField

		results, err := repository.FindPageBy(pageable, parsePathParams(c))
		if err != nil {
			return nil, err
		}

		resources := results.ToResources(HateoasBaseURL(c))
		return &resources, nil
	}, http.StatusPartialContent)
}

// RestFindBy returns all resources matching path params
func RestFindBy(repository repositories.Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (interface{}, error) {
		return repository.FindBy(parsePathParams(c))
	}, http.StatusOK)
}

// RestFindOneBy returns the first resource matching path params
func RestFindOneBy(repository repositories.Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (interface{}, error) {
		result, err := findByPath(c, repository)
		if err != nil {
			return nil, err
		}
		return result, nil
	}, http.StatusOK)
}

// RestRemoveOneBy removes a given resource
func RestRemoveOneBy(repository repositories.Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) error {
		result, err := findByPath(c, repository)
		if err != nil {
			return err
		}
		return repository.Remove(result)
	}, http.StatusNoContent)
}

// RestRemoveAll removes a whole collection
func RestRemoveAll(repository repositories.TruncatableRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) error {
		return repository.Truncate()
	}, http.StatusNoContent)
}

// RestUpsert replace or create a resource
func RestUpsert(repository repositories.Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) error {
		resource := repository.GetNewEntityInstance()
		if err := c.Bind(resource); err != nil {
			return err
		}
		oldres, err := findByPath(c, repository)
		if repositories.IsEntityDoesNotExistError(err) || err == RestErrorGone {
			if c.Request.Method == http.MethodPost {
				c.Header("Location", fmt.Sprintf("%s/%s", c.Request.URL.Path, resource.GetID()))
			}
			if err := repository.Save(resource); err != nil {
				return err
			}
			return RestErrorCreated
		}
		if err != nil {
			return err
		}

		if c.Request.Method == http.MethodPost {
			return errors.AlreadyExistsf("entity '%s' already exists", resource.GetID())
		}
		if err := resource.SetID(oldres.GetID()); err != nil {
			return err
		}
		return repository.Save(resource)
	}, http.StatusOK)
}

// RestErrorHook Convert repository errors in juju errors
func RestErrorHook(tonicErrorHook tonic.ErrorHook) tonic.ErrorHook {
	return func(c *gin.Context, err error) (int, interface{}) {
		if errors.IsAlreadyExists(err) {
			return http.StatusConflict, nil
		}
		switch inner := err.(type) {
		case restErrorCreated:
			return http.StatusCreated, nil
		case restErrorGone:
			return http.StatusGone, nil
		case repositories.EntityDoesNotExistError:
			err = errors.NewNotFound(inner, inner.Error())
		case repositories.UnsupportedIndexError:
			err = errors.NewNotSupported(err, err.Error())
		}
		return tonicErrorHook(c, err)
	}
}

func findByPath(c *gin.Context, repository repositories.Repository) (repositories.Entity, error) {
	params := parsePathParams(c)
	if repo, ok := repository.(repositories.SoftDeletableRepository); ok {
		result, err := repo.FindOneByUnscoped(params)
		if err != nil {
			return nil, err
		}

		if result.GetDeletedAt() != nil {
			return result, RestErrorGone
		}
		return result, err
	}
	return repository.FindOneBy(params)
}

func parsePathParams(c *gin.Context) map[string]interface{} {
	criteria := map[string]interface{}{}
	for _, p := range c.Params {
		criteria[p.Key] = p.Value
	}
	return criteria
}

// HateoasBaseURL returns the base path that has been used to access current resource
func HateoasBaseURL(c *gin.Context) string {
	basePath, ok := c.Get("BasePath")
	if ok {
		return basePath.(string)
	}
	return c.Request.URL.EscapedPath()
}
