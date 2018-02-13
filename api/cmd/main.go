package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/loopfz/gadgeto/tonic/utils/swag"
	"github.com/ovh/lhasa/api/handlers"
)

var version string

// redirect unknown routes to angular
func redirect(c *gin.Context) {
	c.File("./webui/index.html")
}

func main() {

	router := gin.Default()
	// redirect root routes to angular assets
	router.Use(static.Serve("/", static.LocalFile("./webui", true)))

	// redirect unknown routes to angular
	router.NoRoute(redirect)

	// authenticated is the main API route group
	authenticated := router.Group("/api/v0")

	// flushdb flushes the the database for testing purposes
	// TODO: deactivate this route in production
	authenticated.POST("/flushdb", tonic.Handler(handlers.FlushDBHandler, http.StatusOK))

	applicationsRoute := authenticated.Group("/applications")

	applicationsRoute.GET("/", tonic.Handler(handlers.ListApplicationsHandler, http.StatusOK))
	applicationsRoute.POST("/", tonic.Handler(handlers.CreateApplicationHandler, http.StatusCreated))
	applicationsRoute.GET("/:id", tonic.Handler(handlers.DetailApplicationHandler, http.StatusOK))
	applicationsRoute.DELETE("/:id", tonic.Handler(handlers.DeleteApplicationHandler, http.StatusOK))

	// unsecured group does not check incoming signatures
	unsecured := router.Group("/unsecured")
	// health check route
	unsecured.GET("/ping", tonic.Handler(handlers.PingHandler, http.StatusOK))
	// API version
	unsecured.GET("/version", tonic.Handler(handlers.VersionHandler(version), http.StatusOK))
	// auto-generated swagger documentation
	unsecured.GET("/swagger.json", swag.Swagger(router, ""))

	panic(router.Run(":8081"))
}
