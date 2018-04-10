package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	rest "github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1/application"
	"github.com/ovh/lhasa/api/v1/deployment"
	"github.com/ovh/lhasa/api/v1/environment"
)

// registerRoutes registers v1 API routes on a gin engine
func registerRoutes(group *gin.RouterGroup, appRepo *application.Repository, envRepo *environment.Repository, depRepo *deployment.Repository, d deployment.Deployer) {
	appRoutes := group.Group("/applications")
	appRoutes.GET("/", rest.HandlerFindByPage(appRepo))
	appRoutes.DELETE("/", rest.HandlerRemoveAll(appRepo))
	appRoutes.GET("/:domain", rest.HandlerFindByPage(appRepo))
	appRoutes.GET("/:domain/:name", rest.HandlerFindByPage(appRepo))
	appRoutes.GET("/:domain/:name/:version", rest.HandlerFindOneBy(appRepo))
	appRoutes.DELETE("/:domain/:name/:version", rest.HandlerRemoveOneBy(appRepo))
	appRoutes.PUT("/:domain/:name/:version", application.HandlerCreate(appRepo))

	appRoutes.GET("/:domain/:name/:version/deployments/", deployment.HandlerListActiveDeployments(appRepo, depRepo))
	appRoutes.GET("/:domain/:name/:version/deployments/:slug", deployment.HandlerFindDeployment(appRepo, envRepo, depRepo))
	appRoutes.POST("/:domain/:name/:version/deploy/:slug", deployment.HandlerDeploy(appRepo, envRepo, d))

	envRoutes := group.Group("/environments")
	envRoutes.GET("/", rest.HandlerFindByPage(envRepo))
	envRoutes.DELETE("/", rest.HandlerRemoveAll(envRepo))
	envRoutes.GET("/:slug", rest.HandlerFindOneBy(envRepo))
	envRoutes.PUT("/:slug", environment.HandlerCreate(envRepo))
	envRoutes.DELETE("/:slug", rest.HandlerRemoveOneBy(envRepo))

	depRoutes := group.Group("/deployments")
	depRoutes.GET("/", rest.HandlerFindByPage(depRepo))
	depRoutes.DELETE("/", rest.HandlerRemoveAll(depRepo))
	depRoutes.GET("/:public_id", rest.HandlerFindOneBy(depRepo))
	depRoutes.DELETE("/:public_id", rest.HandlerRemoveOneBy(depRepo))
}

// Init initialize the API v1 module
func Init(db *gorm.DB, group *gin.RouterGroup) {
	appRepo := application.NewRepository(db)
	envRepo := environment.NewRepository(db)
	depRepo := deployment.NewRepository(db)
	deployer := deployment.ApplicationDeployer(depRepo)

	registerRoutes(group, appRepo, envRepo, depRepo, deployer)
}
