package routing

import (
	"github.com/jinzhu/gorm"
	"github.com/wI2L/fizz"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
	"github.com/ovh/lhasa/api/v1/application"
	"github.com/ovh/lhasa/api/v1/deployment"
	"github.com/ovh/lhasa/api/v1/environment"
)

// registerRoutes registers v1 API routes on a gin engine
func registerRoutes(group *fizz.RouterGroup, appRepo *application.Repository, envRepo *environment.Repository, depRepo *deployment.Repository, deployer deployment.Deployer, depend deployment.Depend) {
	group.GET("/", []fizz.OperationOption{
		fizz.Summary("Hateoas index of available resources"),
		fizz.ID("IndexV1"),
	}, hateoas.HandlerIndex(
		hateoas.ResourceLink{Href: v1.ApplicationBasePath, Rel: "applications"},
		hateoas.ResourceLink{Href: v1.EnvironmentBasePath, Rel: "environments"},
		hateoas.ResourceLink{Href: v1.DeploymentBasePath, Rel: "deployments"},
	))

	appRoutes := group.Group("/applications", "applications", "Application versions resource management")
	appRoutes.GET("/", getOperationOptions("FindByPage", appRepo,
		fizz.Summary("Find a page of Applications"),
	), hateoas.HandlerFindByPage(appRepo))
	appRoutes.DELETE("/", getOperationOptions("RemoveAll", appRepo,
		fizz.Summary("Delete all Applications"),
	), hateoas.HandlerRemoveAll(appRepo))
	appRoutes.GET("/:domain", getOperationOptions("FindByPageDomain", appRepo,
		fizz.Summary("Find a page of Applications"),
		fizz.InputModel(v1.Application{}),
	), hateoas.HandlerFindByPage(appRepo))
	appRoutes.GET("/:domain/:name", getOperationOptions("FindByPageDomainName", appRepo,
		fizz.Summary("Find a page of Applications"),
		fizz.InputModel(v1.Application{}),
	), hateoas.HandlerFindByPage(appRepo))
	appRoutes.GET("/:domain/:name/:version", getOperationOptions("FindOneBy", appRepo,
		fizz.Summary("Find one Application"),
		fizz.InputModel(v1.Application{}),
	), hateoas.HandlerFindOneBy(appRepo))
	appRoutes.DELETE("/:domain/:name/:version", getOperationOptions("RemoveOneBy", appRepo,
		fizz.Summary("Remove an Application"),
		fizz.InputModel(v1.Application{}),
	), hateoas.HandlerRemoveOneBy(appRepo))
	appRoutes.PUT("/:domain/:name/:version", getOperationOptions("Create", appRepo,
		fizz.Summary("Create an Application Version"),
		fizz.Description("Use this route to create a new application version. The `manifest` field can contains "+
			"any properties useful to track applications in your information system. It is recommended to track it as "+
			"a file in your source-control repository."),
		fizz.StatusDescription("Updated"),
		fizz.Response("201", "Created", nil, nil),
	), application.HandlerCreate(appRepo))

	appRoutes.GET("/:domain/:name/:version/deployments/", getOperationOptions("ListActiveDeployments", appRepo,
		fizz.Summary("List active deployments for this application version"),
		fizz.Description("A deployment is *active* on an environment if it has not been marked as *undeployed*. "+
			"Only a single deployment can be active at a time on a given environment."),
	), deployment.HandlerListActiveDeployments(appRepo, depRepo))
	appRoutes.GET("/:domain/:name/:version/deployments/:slug", getOperationOptions("FindDeployment", appRepo,
		fizz.Summary("Find active deployment for this application version, on this environment"),
	), deployment.HandlerFindDeployment(appRepo, envRepo, depRepo))
	appRoutes.POST("/:domain/:name/:version/deploy/:slug", getOperationOptions("Deploy", appRepo,
		fizz.Summary("Mark this application version as deployed on the given environment"),
		fizz.Description("Note that previous versions of this application on this environments will be marked as undeployed."),
		fizz.Header("location", "URI of the created deployment", nil),
	), deployment.HandlerDeploy(appRepo, envRepo, deployer))

	envRoutes := group.Group("/environments", "environments", "Environments resource management")
	envRoutes.GET("/", getOperationOptions("FindByPage", envRepo,
		fizz.Summary("Find a page of Environments"),
	), hateoas.HandlerFindByPage(envRepo))
	envRoutes.DELETE("/", getOperationOptions("RemoveAll", envRepo,
		fizz.Summary("Delete all Environments"),
	), hateoas.HandlerRemoveAll(envRepo))
	envRoutes.GET("/:slug", getOperationOptions("FindOneBy", envRepo,
		fizz.Summary("Find one Environment"),
		fizz.InputModel(v1.Environment{}),
	), hateoas.HandlerFindOneBy(envRepo))
	envRoutes.PUT("/:slug", getOperationOptions("Create", envRepo,
		fizz.Summary("Create an Environment"),
	), environment.HandlerCreate(envRepo))
	envRoutes.DELETE("/:slug", getOperationOptions("RemoveOneBy", envRepo,
		fizz.Summary("Remove an Environment"),
		fizz.InputModel(v1.Environment{}),
	), hateoas.HandlerRemoveOneBy(envRepo))

	depRoutes := group.Group("/deployments", "deployments", "Deployments resource management")
	depRoutes.GET("/", getOperationOptions("FindByPage", depRepo,
		fizz.Summary("Find a page of Deployments"),
		fizz.InputModel(v1.Deployment{}),
	), hateoas.HandlerFindByPage(depRepo))
	depRoutes.DELETE("/", getOperationOptions("RemoveAll", depRepo,
		fizz.Summary("Delete all Deployments"),
		fizz.InputModel(v1.Deployment{}),
	), hateoas.HandlerRemoveAll(depRepo))
	depRoutes.GET("/:public_id", getOperationOptions("FindOneBy", depRepo,
		fizz.Summary("Find one Deployment"),
		fizz.InputModel(v1.Deployment{}),
	), hateoas.HandlerFindOneBy(depRepo))
	depRoutes.DELETE("/:public_id", getOperationOptions("RemoveOneBy", depRepo,
		fizz.Summary("Remove a Deployment"),
		fizz.InputModel(v1.Deployment{}),
	), hateoas.HandlerRemoveOneBy(depRepo))
	depRoutes.POST("/:public_id/add_link/:target_public_id", getOperationOptions("AddLink", depRepo,
		fizz.Summary("Create a dependency link between two deployments"),
	), deployment.HandlerDepend(depRepo, depend))
}

// Init initialize the API v1 module
func Init(db *gorm.DB, group *fizz.RouterGroup) {
	appRepo := application.NewRepository(db)
	envRepo := environment.NewRepository(db)
	depRepo := deployment.NewRepository(db)
	deployer := deployment.ApplicationDeployer(depRepo)
	depend := deployment.Dependency(depRepo)

	registerRoutes(group, appRepo, envRepo, depRepo, deployer, depend)
}

// getOperationOptions returns an OperationOption list including generated ID for this repository
func getOperationOptions(baseName string, repository hateoas.Repository, options ...fizz.OperationOption) []fizz.OperationOption {
	return append(options, fizz.ID(baseName+repository.GetType().Name()))
}
