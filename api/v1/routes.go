package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/ovh/lhasa/api/handlers"
	app "github.com/ovh/lhasa/api/v1/handlers"
	"github.com/ovh/lhasa/api/v1/repositories"
)

// Register registers v1 API routes on a gin engine
func Register(group *gin.RouterGroup, applicationRepository *repositories.ApplicationRepository) {
	applicationsRoute := group.Group("/applications")
	applicationsRoute.GET("/", tonic.Handler(handlers.RestFindByPage(applicationRepository), http.StatusPartialContent))
	applicationsRoute.DELETE("/", tonic.Handler(handlers.RestRemoveAll(applicationRepository), http.StatusNoContent))
	applicationsRoute.GET("/:domain", tonic.Handler(handlers.RestFindByPage(applicationRepository), http.StatusOK))
	applicationsRoute.GET("/:domain/:name", tonic.Handler(handlers.RestFindByPage(applicationRepository), http.StatusOK))
	applicationsRoute.GET("/:domain/:name/:version", tonic.Handler(handlers.RestFindOneBy(applicationRepository), http.StatusOK))
	applicationsRoute.DELETE("/:domain/:name/:version", tonic.Handler(handlers.RestRemoveOneBy(applicationRepository), http.StatusNoContent))
	applicationsRoute.PUT("/:domain/:name/:version", tonic.Handler(app.ApplicationCreate(applicationRepository), http.StatusOK))
}
