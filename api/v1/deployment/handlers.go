package deployment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/handlers"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
	"github.com/ovh/lhasa/api/v1/application"
	"github.com/ovh/lhasa/api/v1/environment"
)

type deploymentCreateRequest struct {
	*v1.Deployment
	Domain  string `path:"domain" description:"Application Domain"`
	Name    string `path:"name" description:"Application Name"`
	Version string `path:"version" description:"Application Version"`
	Slug    string `path:"slug" description:"Environment identifier"`
}

// HandlerDeploy deploy this application version to the given environment and removes old deployments
func HandlerDeploy(appRepo application.FindOneByUniqueKey, envRepo *environment.Repository, deployer Deployer) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *deploymentCreateRequest) (*v1.Deployment, error) {
		dep := request.Deployment
		app, err := appRepo.FindOneByDomainNameVersion(request.Domain, request.Name, request.Version)
		if err != nil {
			return nil, err
		}
		env, err := envRepo.FindOneBySlug(request.Slug)
		if err != nil {
			return nil, err
		}

		dep, created, err := deployer(*app, *env, dep, handlers.GetLogger(c))
		if err != nil {
			return nil, err
		}
		// If a resource has been created on the origin server, the response SHOULD be 201 (Created) and contain an
		// entity which describes the status of the request and refers to the new resource, and a Location header.
		// https://tools.ietf.org/html/rfc2616#page-54
		if created {
			c.Header("location", dep.GetSelfURL(hateoas.BaseURL(c)))
			c.AbortWithStatusJSON(201, dep)
		}
		return dep, nil
	}, http.StatusCreated)
}

type dependCreateRequest struct {
	PublicID       string `path:"public_id" description:"ID of the deployment that owns the dependency"`
	TargetPublicID string `path:"target_public_id" description:"ID of the deployment targeted by the dependency"`
	Type           string `json:"type"`
}

// HandlerDepend add an observable depdendency with its
func HandlerDepend(depRepo *Repository, depend Depend) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *dependCreateRequest) (hateoas.Entity, error) {
		// Find this dependency by its ID (public)
		entity, err := depRepo.FindOneBy(map[string]interface{}{"public_id": request.PublicID})
		if err != nil {
			return nil, err
		}
		src := entity.(*v1.Deployment)
		target := &v1.Deployment{PublicID: request.TargetPublicID}

		// In depend update dependencies
		if err := depend(src, target, request.Type); err != nil {
			return nil, &(hateoas.InternalError{Message: err.Error(), Detail: request.Type})
		}

		return src, nil
	}, http.StatusOK)
}

// HandlerFindDeployment finds deployment for a given application and environment
func HandlerFindDeployment(appRepo application.FindOneByUniqueKey, envRepo *environment.Repository, depRepo *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *deploymentCreateRequest) (*v1.Deployment, error) {
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
		if dep, ok := res.(*v1.Deployment); ok {
			return dep, nil
		}
		return nil, fmt.Errorf("returned entity %T is not a deployment", res)
	}, http.StatusOK)
}

// HandlerListApplicationActiveDeployments list active deployments for a given application
func HandlerListApplicationActiveDeployments(depRepo *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *deploymentCreateRequest) (interface{}, error) {
		criteria := map[string]interface{}{}
		deps, err := depRepo.FindActivesBy(request.Domain, request.Name, criteria)
		if err != nil {
			return nil, err
		}
		return deps, nil
	}, http.StatusOK)
}

// HandlerListReleaseActiveDeployments list active deployments for a given release (with version)
func HandlerListReleaseActiveDeployments(depRepo *Repository) gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context, request *deploymentCreateRequest) (interface{}, error) {
		criteria := map[string]interface{}{}
		deps, err := depRepo.FindActivesByRelease(request.Domain, request.Name, request.Version, criteria)
		if err != nil {
			return nil, err
		}
		return deps, nil
	}, http.StatusOK)
}
