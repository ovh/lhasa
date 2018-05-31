package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
	"github.com/ovh/lhasa/api/config"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
)

// HandlerCreate replace or create a resource
func HandlerCreate(repository *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, application *v1.ApplicationVersion) error {
		oldres, err := repository.FindOneByUnscoped(map[string]interface{}{"domain": application.Domain, "name": application.Name, "version": application.Version})
		oldapp := oldres.(*v1.ApplicationVersion)
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

type manifestRequest struct {
	Domain     string                 `path:"domain" description:"Application Domain"`
	Name       string                 `path:"name" description:"Application Name"`
	Version    string                 `path:"version" description:"Application Version"`
	Repository string                 `body:"repository" description:"Manifest git assistant"`
	Manifest   map[string]interface{} `body:"manifest" description:"Application Manifest"`
}

// HandlerAssistant help to create manifest with assistant (bitbucket, ...)
func HandlerAssistant(appRepo *Repository, assistant MetaAssistant) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *manifestRequest) (*v1.ApplicationVersion, error) {
		// Find the application
		app, err := appRepo.FindOneByDomainName(request.Domain, request.Name)
		if err != nil {
			return nil, errors.NotFoundf("Application must exist")
		}

		logrus.WithFields(logrus.Fields{
			"request": request,
			"headers": c.GetHeader("X-Remote-User"),
		}).Debug("Pull request assistant")

		if request.Manifest == nil {
			return nil, errors.BadRequestf("Manifest cannot be null")
		}

		if c.GetHeader("X-Remote-User") == "" {
			return nil, errors.Unauthorizedf("User is not authorized")
		}

		name := ""
		manifest, ok := config.ExtractValue("manifest").(map[string]interface{})
		if ok {
			name, _ = manifest["name"].(string)
		}
		parameters := PullRequest{
			Repository:   request.Repository,
			Manifest:     request.Manifest,
			Creator:      c.GetHeader("X-Remote-User"),
			ManifestName: name,
		}

		// launch assistant behaviour
		return app, assistant(*app, &parameters)
	}, http.StatusOK)
}
