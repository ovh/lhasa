package routers

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/loopfz/gadgeto/tonic/utils/jujerr"
	"github.com/loopfz/gadgeto/tonic/utils/swag"
	"github.com/sirupsen/logrus"
	"github.com/ovh/lhasa/api/handlers"
	"github.com/ovh/lhasa/api/v1"
	"github.com/ovh/lhasa/api/v1/repositories"
)

// redirect unknown routes to angular
func redirect(c *gin.Context) {
	c.File("./webui/index.html")
}

//NewRouter creates a new and configured gin router
func NewRouter(applicationRepository *repositories.ApplicationRepository, version string, debugMode bool, log *logrus.Logger) *gin.Engine {
	router := gin.Default()
	configureGin(log, debugMode)

	tonic.SetErrorHook(handlers.RestErrorHook(jujerr.ErrHook))
	// redirect root routes to angular assets
	router.Use(static.Serve("/", static.LocalFile("./webui", true)))

	// redirect unknown routes to angular
	router.NoRoute(redirect)

	v1.Register(router.Group("/api/v1"), applicationRepository)

	// unsecured group does not check incoming signatures
	unsecured := router.Group("/unsecured")
	// health check route
	unsecured.GET("/mon", tonic.Handler(handlers.PingHandler, http.StatusOK))
	// API version
	unsecured.GET("/version", tonic.Handler(handlers.VersionHandler(version), http.StatusOK))
	// auto-generated swagger documentation
	unsecured.GET("/swagger.json", swag.Swagger(router, ""))

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
