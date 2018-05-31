package routers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/loopfz/gadgeto/tonic/utils/jujerr"
	"github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
	"github.com/ovh/lhasa/api/db"
	ext "github.com/ovh/lhasa/api/ext/binding"
	"github.com/ovh/lhasa/api/handlers"
	"github.com/ovh/lhasa/api/hateoas"
	v1 "github.com/ovh/lhasa/api/v1/routing"
)

const uiLocalPath = "./webui"

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
func uiRedirectHandler(uiBasePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		dir, name, notfound := findResource(uiLocalPath+path.Dir(c.Request.URL.Path), path.Base(c.Request.URL.Path))
		for notfound && dir != "." {
			dir, name, notfound = findResource(dir, name)
			if !notfound {
				c.File(dir + name)
				return
			}
		}
		if strings.HasPrefix(c.Request.URL.Path, "/ui") {
			c.Redirect(301, fmt.Sprintf("%s/?redirect=%s", uiBasePath, strings.TrimPrefix(c.Request.URL.Path, "/ui")))
			return
		}
		c.File(uiLocalPath + "/index.html")
	}
}

//NewRouter creates a new and configured gin router
func NewRouter(tm db.TransactionManager, version, hateoasBaseBath, uiBasePath string, debugMode bool, log *logrus.Logger) *fizz.Fizz {
	router := fizz.New()
	router.Generator().OverrideDataType(reflect.TypeOf(&postgres.Jsonb{}), "object", "")

	router.Use(handlers.LoggingMiddleware(log), gin.Recovery())
	configureGin(log, debugMode)

	tonic.SetErrorHook(hateoas.ErrorHook(jujerr.ErrHook))

	// redirect root routes to angular assets
	router.Use(gzip.Gzip(gzip.DefaultCompression), static.Serve("/ui", static.LocalFile("./webui", true)))

	// Set specific hook
	tonic.SetBindHook(ext.BindHook)
	tonic.SetRenderHook(ext.RenderHook, "")

	api := router.Group("/api", "", "", hateoas.AddToBasePath(hateoasBaseBath))
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

	// redirect unknown routes to angular
	router.Engine().NoRoute(uiRedirectHandler(uiBasePath))
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
