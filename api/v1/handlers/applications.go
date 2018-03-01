package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/ovh/lhasa/api/v1/models"
	"github.com/ovh/lhasa/api/v1/repositories"
)

// ApplicationCreate replace or create a resource
func ApplicationCreate(repository *repositories.ApplicationRepository) func(*gin.Context) error {
	return func(c *gin.Context) error {
		application := &models.Application{}
		if err := c.BindJSON(application); err != nil {
			return err
		}
		domain := c.Param("domain")
		name := c.Param("name")
		if application.Domain != domain || application.Name != name {
			return errors.BadRequestf("given application name %s/%s should match with payload %s/%s", application.Domain, application.Name, domain, name)
		}
		application.Version = c.Param("version")
		return repository.Save(application)
	}
}
