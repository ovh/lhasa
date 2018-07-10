package routers

import (
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/loopfz/gadgeto/tonic/utils/jujerr"
	"github.com/ovh/lhasa/api/db"
	ext "github.com/ovh/lhasa/api/ext/binding"
	"github.com/ovh/lhasa/api/handlers"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/security"
	v1 "github.com/ovh/lhasa/api/v1/routing"
	"github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

func uiRedirectHandler(uiBasePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/api/") {
			c.AbortWithStatusJSON(404, gin.H{"error": "API route not found"})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"UIBasePath": uiBasePath})
	}
}

//NewRouter creates a new and configured gin router
func NewRouter(tm db.TransactionManager, version, hateoasBaseBath, uiBasePath string, ServerUIBasePath, webUIDir string, debugMode bool, policy security.Policy, log *logrus.Logger) *fizz.Fizz {
	router := fizz.New()
	router.Generator().OverrideDataType(reflect.TypeOf(&postgres.Jsonb{}), "object", "")

	router.Use(handlers.LoggingMiddleware(log), gin.Recovery())
	configureGin(log, debugMode)

	tonic.SetErrorHook(hateoas.ErrorHook(jujerr.ErrHook))

	// Set specific hook
	tonic.SetBindHook(ext.BindHook)
	tonic.SetRenderHook(ext.RenderHook, "")

	api := router.Group("/api", "", "", hateoas.AddToBasePath(hateoasBaseBath), handlers.AuthMiddleware(policy))
	api.GET("/", []fizz.OperationOption{
		fizz.Summary("Hateoas index of available resources"),
		fizz.ID("IndexAPI"),
	}, hateoas.HandlerIndex(
		hateoas.ResourceLink{Rel: "v1", Href: "/v1"},
		hateoas.ResourceLink{Rel: "unsecured", Href: "/unsecured"},
	))

	v1.Init(tm, api.Group("/v1", "", "", hateoas.AddToBasePath("/v1")), log)

	// unsecured group does not check incoming signatures
	unsecured := api.Group("/unsecured", "unsecured", "Authentication-free routes", hateoas.AddToBasePath("/unsecured"))
	unsecured.GET("/", []fizz.OperationOption{
		fizz.Summary("Hateoas index of available resources"),
		fizz.ID("IndexUnsecured"),
	}, hateoas.HandlerIndex(
		hateoas.ResourceLink{Rel: "monitoring", Href: "/mon"},
		hateoas.ResourceLink{Rel: "version", Href: "/version"},
	))
	// health check route
	unsecured.GET("/mon", []fizz.OperationOption{fizz.ID("Monitoring"), fizz.Summary("Check application and subcomponents health")}, tonic.Handler(handlers.PingHandler(tm.DB().DB()), http.StatusOK))
	// API version
	unsecured.GET("/version", []fizz.OperationOption{fizz.ID("Version"), fizz.Summary("Show the current version of the server")}, tonic.Handler(handlers.VersionHandler(version), http.StatusOK))

	// auto-generated swagger documentation
	infos := &openapi.Info{
		Title: "OpenAPI specification",
		License: &openapi.License{
			Name: "BSD 3-Clause License",
			URL:  "https://opensource.org/licenses/BSD-3-Clause",
		},
		Version: "v1", // this refers to the latest stable api version available on this server
	}
	unsecured.GET("/openapi.json", nil, router.OpenAPI(infos, "json"))
	unsecured.GET("/openapi.yaml", nil, router.OpenAPI(infos, "yaml"))

	if _, err := os.Stat(webUIDir + "/index.html"); os.IsNotExist(err) {
		log.Warn("index.html not found. Starting in API only mode. No static content will be served.")
		return router
	}
	// serve static content from angular
	router.Use(static.Serve(ServerUIBasePath, static.LocalFile(webUIDir, false)))

	// redirect unknown routes to angular
	router.Engine().LoadHTMLFiles(webUIDir + "/index.html")
	router.Engine().NoRoute(gzip.Gzip(gzip.DefaultCompression), uiRedirectHandler(uiBasePath))
	return router
}

func configureGin(log *logrus.Logger, debugMode bool) {
	if log != nil {
		gin.DefaultWriter = log.Writer()
	}
	ginMode := gin.ReleaseMode
	if debugMode {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)
}
