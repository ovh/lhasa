package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	rest "github.com/ovh/lhasa/api/handlers"
	v1 "github.com/ovh/lhasa/api/v1/handlers"
	"github.com/ovh/lhasa/api/v1/repositories"
	"github.com/ovh/lhasa/api/v1/service"
)

// registerRoutes registers v1 API routes on a gin engine
func registerRoutes(group *gin.RouterGroup, appRepo *repositories.ApplicationRepository, envRepo *repositories.EnvironmentRepository, depRepo *repositories.DeploymentRepository, deployer service.Deployer) {
	appRoutes := group.Group("/applications")
	appRoutes.GET("/", rest.RestFindByPage(appRepo))
	appRoutes.DELETE("/", rest.RestRemoveAll(appRepo))
	appRoutes.GET("/:domain", rest.RestFindByPage(appRepo))
	appRoutes.GET("/:domain/:name", rest.RestFindByPage(appRepo))
	appRoutes.GET("/:domain/:name/:version", rest.RestFindOneBy(appRepo))
	appRoutes.DELETE("/:domain/:name/:version", rest.RestRemoveOneBy(appRepo))
	appRoutes.PUT("/:domain/:name/:version", v1.ApplicationCreate(appRepo))

	appRoutes.GET("/:domain/:name/:version/deployments/", v1.ApplicationListActiveDeployments(appRepo, depRepo))
	appRoutes.GET("/:domain/:name/:version/deployments/:slug", v1.ApplicationFindLastDeployment(appRepo, envRepo, depRepo))
	appRoutes.POST("/:domain/:name/:version/deploy/:slug", v1.ApplicationDeploy(appRepo, envRepo, deployer))

	envRoutes := group.Group("/environments")
	envRoutes.GET("/", rest.RestFindByPage(envRepo))
	envRoutes.DELETE("/", rest.RestRemoveAll(envRepo))
	envRoutes.GET("/:slug", rest.RestFindOneBy(envRepo))
	envRoutes.PUT("/:slug", v1.EnvironmentCreate(envRepo))
	envRoutes.DELETE("/:slug", rest.RestRemoveOneBy(envRepo))

	depRoutes := group.Group("/deployments")
	depRoutes.GET("/", rest.RestFindByPage(depRepo))
	depRoutes.DELETE("/", rest.RestRemoveAll(depRepo))
	depRoutes.GET("/:public_id", rest.RestFindOneBy(depRepo))
	depRoutes.DELETE("/:public_id", rest.RestRemoveOneBy(depRepo))
}

// Init initialize the API v1 module
func Init(db *gorm.DB, group *gin.RouterGroup) {
	appRepo := repositories.NewApplicationRepository(db)
	envRepo := repositories.NewEnvironmentRepository(db)
	depRepo := repositories.NewDeploymentRepository(db)
	deployer := service.ApplicationDeployer(db, depRepo)

	registerRoutes(group, appRepo, envRepo, depRepo, deployer)
}
