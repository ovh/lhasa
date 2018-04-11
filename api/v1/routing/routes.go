package routing

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
	"github.com/ovh/lhasa/api/v1/application"
	"github.com/ovh/lhasa/api/v1/deployment"
	"github.com/ovh/lhasa/api/v1/environment"
)

// registerRoutes registers v1 API routes on a gin engine
func registerRoutes(group *gin.RouterGroup, appRepo *application.Repository, envRepo *environment.Repository, depRepo *deployment.Repository, d deployment.Deployer) {
	group.GET("/", indexHandler())

	appRoutes := group.Group("/applications")
	appRoutes.GET("/", hateoas.HandlerFindByPage(appRepo))
	appRoutes.DELETE("/", hateoas.HandlerRemoveAll(appRepo))
	appRoutes.GET("/:domain", hateoas.HandlerFindByPage(appRepo))
	appRoutes.GET("/:domain/:name", hateoas.HandlerFindByPage(appRepo))
	appRoutes.GET("/:domain/:name/:version", hateoas.HandlerFindOneBy(appRepo))
	appRoutes.DELETE("/:domain/:name/:version", hateoas.HandlerRemoveOneBy(appRepo))
	appRoutes.PUT("/:domain/:name/:version", application.HandlerCreate(appRepo))

	appRoutes.GET("/:domain/:name/:version/deployments/", deployment.HandlerListActiveDeployments(appRepo, depRepo))
	appRoutes.GET("/:domain/:name/:version/deployments/:slug", deployment.HandlerFindDeployment(appRepo, envRepo, depRepo))
	appRoutes.POST("/:domain/:name/:version/deploy/:slug", deployment.HandlerDeploy(appRepo, envRepo, d))

	envRoutes := group.Group("/environments")
	envRoutes.GET("/", hateoas.HandlerFindByPage(envRepo))
	envRoutes.DELETE("/", hateoas.HandlerRemoveAll(envRepo))
	envRoutes.GET("/:slug", hateoas.HandlerFindOneBy(envRepo))
	envRoutes.PUT("/:slug", environment.HandlerCreate(envRepo))
	envRoutes.DELETE("/:slug", hateoas.HandlerRemoveOneBy(envRepo))

	depRoutes := group.Group("/deployments")
	depRoutes.GET("/", hateoas.HandlerFindByPage(depRepo))
	depRoutes.DELETE("/", hateoas.HandlerRemoveAll(depRepo))
	depRoutes.GET("/:public_id", hateoas.HandlerFindOneBy(depRepo))
	depRoutes.DELETE("/:public_id", hateoas.HandlerRemoveOneBy(depRepo))
}

func indexHandler() gin.HandlerFunc {
	return tonic.Handler(func(c *gin.Context) (IndexResource, error) {
		return IndexResource{
			Resource: hateoas.Resource{
				Links: []hateoas.ResourceLink{
					{Href: v1.ApplicationBasePath, Rel: "applications"},
					{Href: v1.EnvironmentBasePath, Rel: "environments"},
					{Href: v1.DeploymentBasePath, Rel: "deployments"},
				},
			},
		}, nil
	}, http.StatusOK)
}

type IndexResource struct {
	hateoas.Resource
}

// Init initialize the API v1 module
func Init(db *gorm.DB, group *gin.RouterGroup) {
	appRepo := application.NewRepository(db)
	envRepo := environment.NewRepository(db)
	depRepo := deployment.NewRepository(db)
	deployer := deployment.ApplicationDeployer(depRepo)

	registerRoutes(group, appRepo, envRepo, depRepo, deployer)
}
