package routers

import (
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/loopfz/gadgeto/tonic/utils/jujerr"
	"github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
	"github.com/ovh/lhasa/api/handlers"
	"github.com/ovh/lhasa/api/hateoas"
	v1 "github.com/ovh/lhasa/api/v1/routing"
)

// checkHTML5Path check for HTML5 path in request
func checkHTML5Path(c *gin.Context) bool {
	if strings.HasPrefix(c.Request.URL.Path, "/api") {
		return false
	}
	return true
}

// find existing resource on base webui (for security issue)
func findResource(dir string, name string) (string, string, bool) {
	var existing = strings.Replace(dir+"/"+name, "//", "", -1)
	if _, err := os.Stat(existing); err == nil {
		return "", existing, false
	}
	if dir == "/" {
		return "", name, true
	}
	return path.Dir(dir), name, true
}

// redirect unknown routes to angular
func redirect(c *gin.Context) {
	var basepath = "./webui"
	if checkHTML5Path(c) {
		dir, name, notfound := findResource(basepath+path.Dir(c.Request.URL.Path), path.Base(c.Request.URL.Path))
		for notfound && dir != "." {
			dir, name, notfound = findResource(dir, name)
			if !notfound {
				c.File(dir + name)
				return
			}
		}
	}
	if len(c.Request.URL.Path) > 1 {
		// Path is not slash
		c.Redirect(301, "/?redirect="+c.Request.URL.Path)
		return
	}
	c.File(basepath + "/index.html")
}

//NewRouter creates a new and configured gin router
func NewRouter(db *gorm.DB, version, hateoasBaseBath string, debugMode bool, log *logrus.Logger) *fizz.Fizz {
	router := fizz.New()
	router.Generator().OverrideDataType(reflect.TypeOf(&postgres.Jsonb{}), "object", "")

	router.Use(handlers.LoggingMiddleware(log), gin.Recovery())
	configureGin(log, debugMode)

	tonic.SetErrorHook(hateoas.ErrorHook(jujerr.ErrHook))

	// redirect root routes to angular assets
	router.Use(gzip.Gzip(gzip.DefaultCompression), static.Serve("/", static.LocalFile("./webui", true)))

	api := router.Group("/api", "", "", hateoas.AddToBasePath(hateoasBaseBath))
	api.GET("/", []fizz.OperationOption{
		fizz.Summary("Hateoas index of available resources"),
		fizz.ID("IndexAPI"),
	}, hateoas.HandlerIndex(
		hateoas.ResourceLink{Rel: "v1", Href: "/v1"},
		hateoas.ResourceLink{Rel: "unsecured", Href: "/unsecured"},
	))

	v1.Init(db, api.Group("/v1", "", "", hateoas.AddToBasePath("/v1")))

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
	unsecured.GET("/mon", []fizz.OperationOption{fizz.ID("Monitoring"), fizz.Summary("Check application and subcomponents health")}, tonic.Handler(handlers.PingHandler(db.DB()), http.StatusOK))
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

	// redirect unknown routes to angular
	router.Engine().NoRoute(redirect)
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
