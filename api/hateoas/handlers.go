package hateoas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
)

// InternalError all internal error
type InternalError struct {
	Message string
	Detail  string
}

type errorCreated string
type errorGone string

// Error implements error interface
func (err *InternalError) Error() string {
	return err.Message + ":" + err.Detail
}

// Error implements error interface
func (err errorCreated) Error() string {
	return string(err)
}

// Error implements error interface
func (err errorGone) Error() string {
	return string(err)
}

// ErrorCreated is raised when no error occurs but a resource has been created (tonic single-status code workaround)
var ErrorCreated = errorCreated("created")

// ErrorGone is raised when a former resource has been requested but no longer exist
var ErrorGone = errorGone("gone")

// HandlerFindByPage returns a filtered and paginated resource list
func HandlerFindByPage(repository PageableRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (*PagedResources, error) {
		pageable := Pageable{}
		c.ShouldBindQuery(&pageable)

		// params and query are user to filter resultset
		results, err := repository.FindPageBy(pageable, parsePathParamsAndQuery(c))
		if err != nil {
			return nil, err
		}

		resources := results.ToResources(BaseURL(c))
		return &resources, nil
	}, http.StatusPartialContent)
}

// HandlerFindBy returns all resources matching path params
func HandlerFindBy(repository Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (interface{}, error) {
		return repository.FindBy(parsePathParams(c))
	}, http.StatusOK)
}

// HandlerFindOneBy returns the first resource matching path params
func HandlerFindOneBy(repository Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (interface{}, error) {
		result, err := findByPath(c, repository)
		if err != nil {
			return nil, err
		}
		return result, nil
	}, http.StatusOK)
}

// HandlerRemoveOneBy removes a given resource
func HandlerRemoveOneBy(repository Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) error {
		result, err := findByPath(c, repository)
		if err != nil {
			return err
		}
		return repository.Remove(result)
	}, http.StatusNoContent)
}

// HandlerRemoveAll removes a whole collection
func HandlerRemoveAll(repository TruncatableRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) error {
		return repository.Truncate()
	}, http.StatusNoContent)
}

// HandlerUpsert replace or create a resource
func HandlerUpsert(repository Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) error {
		resource := repository.GetNewEntityInstance()
		if err := c.Bind(resource); err != nil {
			return err
		}
		oldres, err := findByPath(c, repository)
		if IsEntityDoesNotExistError(err) || err == ErrorGone {
			if c.Request.Method == http.MethodPost {
				c.Header("Location", fmt.Sprintf("%s/%s", c.Request.URL.Path, resource.GetID()))
			}
			if err := repository.Save(resource); err != nil {
				return err
			}
			return ErrorCreated
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

// ErrorHook Convert repository errors in juju errors
func ErrorHook(tonicErrorHook tonic.ErrorHook) tonic.ErrorHook {
	return func(c *gin.Context, err error) (int, interface{}) {
		if errors.IsAlreadyExists(err) {
			return http.StatusConflict, nil
		}
		switch inner := err.(type) {
		case errorCreated:
			return http.StatusCreated, nil
		case errorGone:
			return http.StatusGone, nil
		case EntityDoesNotExistError:
			err = errors.NewNotFound(inner, inner.Error())
		case UnsupportedIndexError:
			err = errors.NewNotSupported(err, err.Error())
		}
		return tonicErrorHook(c, err)
	}
}

func findByPath(c *gin.Context, repository Repository) (Entity, error) {
	params := parsePathParams(c)
	if repo, ok := repository.(SoftDeletableRepository); ok {
		result, err := repo.FindOneByUnscoped(params)
		if err != nil {
			return nil, err
		}

		if result.GetDeletedAt() != nil {
			return result, ErrorGone
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

// unCamelCase discover physical column name in database
func unCamelCase(value string) string {
	if strings.Contains(value, "->>") {
		// on jsonb column don't un camel case field
		return value
	}
	accu := ""
	for i := 0; i < len(value); i = i + 1 {
		switch value[i] {
		case '_':
			break
		default:
			// is already lower
			if bytes.ToLower([]byte{value[i]})[0] == value[i] {
				accu = accu + string(value[i])
			} else {
				accu = accu + "_" + string(bytes.ToLower([]byte{value[i]})[0])
			}
		}
	}
	return accu
}

// parse param and query of request
func parsePathParamsAndQuery(c *gin.Context) map[string]interface{} {
	criteria := map[string]interface{}{}
	for _, p := range c.Params {
		criteria[p.Key] = p.Value
	}
	pageable := Pageable{}
	c.ShouldBindQuery(&pageable)

	// scan for optionnal query
	var q = map[string]interface{}{}
	if len(pageable.Query) > 0 {
		json.Unmarshal([]byte(pageable.Query), &q)
	}

	// Scan it
	for k, v := range q {
		if len(k) > 0 {
			// only support simple fields, no check on jsonb expression
			if !strings.Contains(k, ".") {
				// Only support taking first occurence
				str, check := v.(string)
				if !check {
					value, _ := json.Marshal(v)
					str = string(value)
				}
				criteria[unCamelCase(k)] = str
			} else {
				// Only support taking first occurence
				criteria[k] = v.(string)
			}
		}
	}
	return criteria
}

// BaseURL returns the base path that has been used to access current resource
func BaseURL(c *gin.Context) string {
	basePath, ok := c.Get("BasePath")
	if ok {
		return basePath.(string)
	}
	return c.Request.URL.EscapedPath()
}
