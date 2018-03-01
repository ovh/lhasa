package v0

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/v0/handlers"
)

//Register registers v0 API routes on a gin engine
func Register(group *gin.RouterGroup) {
	// flushdb flushes the the database for testing purposes
	group.POST("/flushdb", tonic.Handler(handlers.FlushDBHandler, http.StatusOK))
	// TODO: deactivate this route in production
	applicationsRoute := group.Group("/applications")
	applicationsRoute.GET("/", tonic.Handler(handlers.ListApplicationsHandler, http.StatusOK))
	applicationsRoute.POST("/", tonic.Handler(handlers.CreateApplicationHandler, http.StatusCreated))
	applicationsRoute.GET("/:id", tonic.Handler(handlers.DetailApplicationHandler, http.StatusOK))
	applicationsRoute.DELETE("/:id", tonic.Handler(handlers.DeleteApplicationHandler, http.StatusOK))
}
