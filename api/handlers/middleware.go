package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ovh/lhasa/api/hateoas"
)

var (
	headerSource    = http.CanonicalHeaderKey("X-Ovh-Gateway-Source")
	requestIDHeader = http.CanonicalHeaderKey("X-Request-Id")
)

// LoggingMiddleware logs before and after incoming gin requests
func LoggingMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if log == nil {
			c.Next()
			return
		}
		fields := logrus.Fields{
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"source_token": c.Request.Header.Get(headerSource),
			"request_id":   c.Request.Header.Get(requestIDHeader),
		}
		log.WithFields(fields).Debug("incoming request")
		startTime := time.Now()
		c.Set("logfields", fields)

		c.Next()

		log.WithFields(fields).WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"duration": time.Since(startTime).Seconds(),
		}).Info("done")

		for _, err := range c.Errors.Errors() {
			if err != hateoas.ErrorCreated.Error() {
				log.WithFields(fields).Error(err)
			}
		}
	}
}

// AddToBasePath add a subpath to the BasePath stored in the gin context
func AddToBasePath(basePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path, ok := c.Get("BasePath")
		if !ok {
			c.Set("BasePath", basePath)
			c.Next()
			return
		}
		c.Set("BasePath", path.(string)+basePath)
	}
}
