package routing

import (
	"github.com/ovh/lhasa/api/db"
	"github.com/ovh/lhasa/api/graphapi"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/v1"
	"github.com/ovh/lhasa/api/v1/application"
	"github.com/ovh/lhasa/api/v1/badge"
	"github.com/ovh/lhasa/api/v1/content"
	"github.com/ovh/lhasa/api/v1/deployment"
	"github.com/ovh/lhasa/api/v1/domain"
	"github.com/ovh/lhasa/api/v1/environment"
	"github.com/ovh/lhasa/api/v1/graph"
	"github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
)

// registerRoutes registers v1 API routes on a gin engine
func registerRoutes(group *fizz.RouterGroup,
	graphRepo *graph.Repository,
	domRepo *domain.Repository,
	appRepo *application.Repository,
	appLatestRepo *application.RepositoryLatest,
	contRepo *content.Repository,
	envRepo *environment.Repository,
	depRepo *deployment.Repository,
	badgeRepo *badge.Repository,
	deployer deployment.Deployer,
	depend deployment.Depend,
	latestUpdater application.LatestUpdater) {
	group.GET("/", []fizz.OperationOption{
		fizz.Summary("Hateoas index of available resources"),
		fizz.ID("IndexV1"),
	}, hateoas.HandlerIndex(
		hateoas.ResourceLink{Href: v1.ContentBasePath, Rel: "contents"},
		hateoas.ResourceLink{Href: v1.ApplicationBasePath, Rel: "applications"},
		hateoas.ResourceLink{Href: v1.EnvironmentBasePath, Rel: "environments"},
		hateoas.ResourceLink{Href: v1.DeploymentBasePath, Rel: "deployments"},
	))

	graphRoutes := group.Group("/graphs", "graph", "Graphs node and edge management")
	graphRoutes.GET("/", getGraphOperationOptions("FindAllActive", graphRepo,
		fizz.Summary("Find a page of node and all associated edge"),
	), graphapi.HandlerFindAllActive(graphRepo))

	domRoutes := group.Group("/domains", "domains", "Domains resource management")
	domRoutes.GET("/", getOperationOptions("FindByPage", domRepo,
		fizz.Summary("Find a page of Domains"),
	), hateoas.HandlerFindByPage(domRepo))
	domRoutes.GET("/:domain", getOperationOptions("FindOneBy", domRepo,
		fizz.Summary("Find one Domain"),
		fizz.InputModel(v1.Domain{}),
	), hateoas.HandlerFindOneBy(domRepo))

	contRoutes := group.Group("/contents", "contents", "Content resource management")
	contRoutes.DELETE("/", getOperationOptions("RemoveAll", contRepo,
		fizz.Summary("Delete all Contents"),
	), hateoas.HandlerRemoveAll(contRepo))
	contRoutes.GET("/:name", getOperationOptions("FindOneByName", contRepo,
		fizz.Summary("Find one Content"),
		fizz.InputModel(v1.Content{}),
	), content.HandlerGet(contRepo))
	contRoutes.GET("/:name/:locale", getOperationOptions("FindOneByNameAndLocale", contRepo,
		fizz.Summary("Find one Content"),
		fizz.InputModel(v1.Content{}),
	), content.HandlerGet(contRepo))
	contRoutes.DELETE("/:name", getOperationOptions("RemoveOneBy", contRepo,
		fizz.Summary("Remove an Content"),
		fizz.InputModel(v1.Content{}),
	), hateoas.HandlerRemoveOneBy(contRepo))
	contRoutes.PUT("/:name/:locale", getOperationOptions("Create", contRepo,
		fizz.Summary("Create an Content"),
		fizz.Description("Use this route to create a new content. The `body` field must be plain raw text."),
		fizz.StatusDescription("Updated"),
		fizz.Response("201", "Created", nil, nil),
	), content.HandlerCreate(contRepo))

	appRoutes := group.Group("/applications", "applications", "Application management")
	appRoutes.GET("/", getOperationOptions("FindByPage", appRepo,
		fizz.Summary("Find a page of Applications"),
	), hateoas.HandlerFindByPage(appLatestRepo))
	appRoutes.DELETE("/", getOperationOptions("RemoveAll", appRepo,
		fizz.Summary("Delete all Applications"),
	), hateoas.HandlerRemoveAll(appRepo))
	appRoutes.GET("/:domain", getOperationOptions("FindLatestByPageDomain", appLatestRepo,
		fizz.Summary("Find a page of Applications"),
		fizz.InputModel(v1.Domain{}),
	), hateoas.HandlerFindByPage(appLatestRepo))
	appRoutes.GET("/:domain/:name/latest", getOperationOptions("FindLatestBy", appLatestRepo,
		fizz.Summary("Find latest available version for an Application"),
	), application.HandlerRedirectLatest(appLatestRepo))
	appRoutes.GET("/:domain/:name/versions", getOperationOptions("FindPageOfReleases", appRepo,
		fizz.Summary("Find a page of Releases"),
		fizz.InputModel(v1.Application{}),
	), hateoas.HandlerFindByPage(appRepo))
	appRoutes.GET("/:domain/:name/deployments", getOperationOptions("ListActiveDeploymentsLatest", appLatestRepo,
		fizz.Summary("List active deployments for the latest release of this application"),
		fizz.Description("A deployment is *active* on an environment if it has not been marked as *undeployed*. "+
			"Only a single deployment can be active at a time on a given environment."),
	), deployment.HandlerListApplicationActiveDeployments(appLatestRepo, depRepo))
	appRoutes.GET("/:domain/:name/versions/:version", getOperationOptions("FindOneBy", appRepo,
		fizz.Summary("Find one Release"),
		fizz.InputModel(v1.Release{}),
	), hateoas.HandlerFindOneBy(appRepo))
	appRoutes.DELETE("/:domain/:name/versions/:version", getOperationOptions("RemoveOneBy", appRepo,
		fizz.Summary("Remove a Release"),
		fizz.InputModel(v1.Release{}),
	), hateoas.HandlerRemoveOneBy(appRepo))
	appRoutes.PUT("/:domain/:name/versions/:version", getOperationOptions("Create", appRepo,
		fizz.Summary("Create a Release"),
		fizz.Description("Use this route to create a new application version. The `manifest` field can contains "+
			"any properties useful to track applications in your information system. It is recommended to track it as "+
			"a file in your source-control repository."),
		fizz.StatusDescription("Updated"),
		fizz.Response("201", "Created", nil, nil),
	), application.HandlerCreate(appRepo, latestUpdater))

	appRoutes.GET("/:domain/:name/versions/:version/deployments", getOperationOptions("ListActiveDeploymentsVersion", appRepo,
		fizz.Summary("List active deployments for this application version"),
		fizz.Description("A deployment is *active* on an environment if it has not been marked as *undeployed*. "+
			"Only a single deployment can be active at a time on a given environment."),
	), deployment.HandlerListReleaseActiveDeployments(appRepo, depRepo))
	appRoutes.GET("/:domain/:name/versions/:version/deployments/:slug", getOperationOptions("FindDeployment", appRepo,
		fizz.Summary("Find active deployment for this application version, on this environment"),
	), deployment.HandlerFindDeployment(appRepo, envRepo, depRepo))
	appRoutes.POST("/:domain/:name/versions/:version/deploy/:slug", getOperationOptions("Deploy", appRepo,
		fizz.Summary("Mark this application version as deployed on the given environment"),
		fizz.Description("Note that previous versions of this application on this environments will be marked as undeployed."),
		fizz.Header("location", "URI of the created deployment", nil),
	), deployment.HandlerDeploy(appRepo, envRepo, deployer))
	appRoutes.GET("/:domain/:name/versions/:version/badges", getOperationOptions("FindBadgeRatingsForAnApplicationVersion", appRepo,
		fizz.Summary("Find badge values for an application version"),
	), application.HandlerGetBadgeRatingsForAppVersion(appRepo))
	appRoutes.PUT("/:domain/:name/versions/:version/badgeratings/:badgeslug", getOperationOptions("SetBadgeRatingForAnApplicationVersion", appRepo,
		fizz.Summary("Set badge value for an application version"),
	), application.HandlerSetBadgeRatingForAppVersion(appRepo))
	appRoutes.DELETE("/:domain/:name/versions/:version/badgeratings/:badgeslug", getOperationOptions("DeleteBadgeRatingForAnApplicationVersion", appRepo,
		fizz.Summary("Delete badge value for an application version"),
	), application.HandlerDeleteBadgeRatingForAppVersion(appRepo))

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

	badgeRoutes := group.Group("/badges", "badges", "Badges resource management")
	badgeRoutes.GET("/", getOperationOptions("FindByPage", badgeRepo,
		fizz.Summary("Find a page of Badges"),
	), hateoas.HandlerFindByPage(badgeRepo))
	badgeRoutes.DELETE("/", getOperationOptions("RemoveAll", badgeRepo,
		fizz.Summary("Delete all Badges"),
	), hateoas.HandlerRemoveAll(badgeRepo))
	badgeRoutes.GET("/:slug", getOperationOptions("FindOneBy", badgeRepo,
		fizz.Summary("Find one Badge"),
		fizz.InputModel(v1.Badge{}),
	), hateoas.HandlerFindOneBy(badgeRepo))
	badgeRoutes.PUT("/:slug", getOperationOptions("Create", badgeRepo,
		fizz.Summary("Create a Badge"),
	), badge.HandlerCreate(badgeRepo))
	badgeRoutes.DELETE("/:slug", getOperationOptions("RemoveOneBy", badgeRepo,
		fizz.Summary("Remove a Badge"),
		fizz.InputModel(v1.Badge{}),
	), hateoas.HandlerRemoveOneBy(badgeRepo))

}

// Init initialize the API v1 module
func Init(tm db.TransactionManager, group *fizz.RouterGroup, log *logrus.Logger) {
	graphRepo := graph.NewRepository(tm.DB())
	domRepo := domain.NewRepository(tm.DB())
	appRepo := application.NewRepository(tm.DB())
	appLatestRepo := application.NewRepositoryLatest(tm.DB())
	contRepo := content.NewRepository(tm.DB())
	envRepo := environment.NewRepository(tm.DB())
	badgeRepo := badge.NewRepository(tm.DB())
	deployer := deployment.ApplicationDeployer(tm, deployment.NewRepository)
	depRepo := deployment.NewRepository(tm.DB())
	depend := deployment.Dependency(depRepo)
	latestUpdater := application.NewLatestUpdater(tm, application.NewRepository, application.NewRepositoryLatest, log)

	registerRoutes(group, graphRepo, domRepo, appRepo, appLatestRepo, contRepo, envRepo, depRepo, badgeRepo, deployer, depend, latestUpdater)
}

// getOperationOptions returns an OperationOption list including generated ID for this repository
func getOperationOptions(baseName string, repository hateoas.Repository, options ...fizz.OperationOption) []fizz.OperationOption {
	return append(options, fizz.ID(baseName+repository.GetType().Name()))
}

// getOperationOptions returns an OperationOption list including generated ID for this repository
func getGraphOperationOptions(baseName string, repository graphapi.Repository, options ...fizz.OperationOption) []fizz.OperationOption {
	return append(options, fizz.ID(baseName+repository.GetType().Name()))
}
