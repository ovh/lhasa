package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/handlers"
	"github.com/ovh/lhasa/api/repositories"
	"github.com/ovh/lhasa/api/v1/models"
	v1repo "github.com/ovh/lhasa/api/v1/repositories"
	"github.com/ovh/lhasa/api/v1/service"
)

type deploymentCreateRequest struct {
	*models.Deployment
	Domain  string `path:"domain"`
	Name    string `path:"name"`
	Version string `path:"version"`
	Slug    string `path:"slug"`
}

// ApplicationCreate replace or create a resource
func ApplicationCreate(repository *v1repo.ApplicationRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, application *models.Application) error {
		oldres, err := repository.FindOneByUnscoped(map[string]interface{}{"domain": application.Domain, "name": application.Name, "version": application.Version})
		oldapp := oldres.(*models.Application)
		if repositories.IsEntityDoesNotExistError(err) {
			if err := repository.Save(application); err != nil {
				return err
			}
			return handlers.RestErrorCreated
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
			return handlers.RestErrorCreated
		}
		return nil

	}, http.StatusOK)
}

// ApplicationDeploy deploy this application version to the given environment and removes old deployments
func ApplicationDeploy(appRepo *v1repo.ApplicationRepository, envRepo *v1repo.EnvironmentRepository, deployer service.Deployer) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *deploymentCreateRequest) error {
		dep := request.Deployment
		app, err := appRepo.FindOneByDomainNameVersion(request.Domain, request.Name, request.Version)
		if err != nil {
			return err
		}
		env, err := envRepo.FindOneBySlug(request.Slug)
		if err != nil {
			return err
		}

		if err := deployer(*app, *env, dep); err != nil {
			return err
		}

		return handlers.RestErrorCreated
	}, http.StatusOK)
}

// ApplicationFindLastDeployment list active deployments for a given application
func ApplicationFindLastDeployment(appRepo *v1repo.ApplicationRepository, envRepo *v1repo.EnvironmentRepository, depRepo *v1repo.DeploymentRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *deploymentCreateRequest) (*models.Deployment, error) {
		app, err := appRepo.FindOneByDomainNameVersion(request.Domain, request.Name, request.Version)
		if err != nil {
			return nil, err
		}
		env, err := envRepo.FindOneBySlug(request.Slug)
		if err != nil {
			return nil, err
		}
		res, err := depRepo.FindOneBy(map[string]interface{}{"application_id": app.ID, "environment_id": env.ID})
		if err != nil {
			return nil, err
		}
		dep := res.(*models.Deployment)
		dep.ToResource(handlers.HateoasBaseURL(c))
		if err != nil {
			return nil, err
		}
		return dep, nil
	}, http.StatusOK)
}

// ApplicationListActiveDeployments list active deployments for a given application
func ApplicationListActiveDeployments(appRepo *v1repo.ApplicationRepository, depRepo *v1repo.DeploymentRepository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *deploymentCreateRequest) (interface{}, error) {
		app, err := appRepo.FindOneByDomainNameVersion(request.Domain, request.Name, request.Version)
		if err != nil {
			return nil, err
		}
		criteria := map[string]interface{}{"application_id": app.ID}
		deps, err := depRepo.FindActivesBy(criteria)
		if err != nil {
			return nil, err
		}
		for _, dep := range deps {
			dep.ToResource(handlers.HateoasBaseURL(c))
		}
		return deps, nil
	}, http.StatusOK)
}
