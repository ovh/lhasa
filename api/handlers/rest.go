package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/repositories"
)

// RestFindByPage returns a filtered and paginated resource list
func RestFindByPage(repository repositories.PageableRepository) func(c *gin.Context) (*repositories.PagedResources, error) {
	return func(c *gin.Context) (*repositories.PagedResources, error) {
		pageable := repositories.Pageable{}
		c.ShouldBindQuery(&pageable)

		results, err := repository.FindPageBy(pageable, parsePathParams(c))
		if err != nil {
			return nil, err
		}

		resources := results.ToResources(hateoasBaseURL(c))
		return &resources, nil
	}
}

// RestFindByIndexedPage returns a filtered and paginated resource list
func RestFindByIndexedPage(repository repositories.PageableRepository, indexField string) func(c *gin.Context) (*repositories.PagedResources, error) {
	return func(c *gin.Context) (*repositories.PagedResources, error) {
		pageable := repositories.Pageable{}
		c.ShouldBindQuery(&pageable)
		pageable.IndexedBy = indexField

		results, err := repository.FindPageBy(pageable, parsePathParams(c))
		if err != nil {
			return nil, err
		}

		resources := results.ToResources(hateoasBaseURL(c))
		return &resources, nil
	}
}

// RestFindBy returns all resources matching path params
func RestFindBy(repository repositories.Repository) func(*gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		return repository.FindBy(parsePathParams(c))
	}
}

// RestFindOneBy returns the first resource matching path params
func RestFindOneBy(repository repositories.Repository) func(*gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		result, err := repository.FindOneBy(parsePathParams(c))
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

// RestRemoveOneBy removes a given resource
func RestRemoveOneBy(repository repositories.Repository) func(*gin.Context) error {
	return func(c *gin.Context) error {
		result, err := repository.FindOneBy(parsePathParams(c))
		if err != nil {
			return err
		}
		return repository.Remove(result)
	}
}

// RestCreate replace or create a resource
func RestCreate(repository repositories.Repository, resource interface{}) func(*gin.Context) error {
	return func(c *gin.Context) error {
		if err := c.BindJSON(resource); err != nil {
			return err
		}
		return repository.Save(resource)
	}
}

// RestErrorHook Convert repository errors in juju errors
func RestErrorHook(tonicErrorHook tonic.ErrorHook) tonic.ErrorHook {
	return func(c *gin.Context, err error) (int, interface{}) {
		switch inner := err.(type) {
		case repositories.EntityDoesNotExistError:
			err = errors.NewNotFound(inner, inner.Error())
		case repositories.UnsupportedIndexError:
			err = errors.NewNotSupported(err, err.Error())
		}
		return tonicErrorHook(c, err)
	}
}

func parsePathParams(c *gin.Context) map[string]interface{} {
	criteria := map[string]interface{}{}
	for _, p := range c.Params {
		criteria[p.Key] = p.Value
	}
	return criteria
}

func hateoasBaseURL(c *gin.Context) string {
	return c.Request.URL.EscapedPath()
}
