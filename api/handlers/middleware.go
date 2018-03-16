package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	headerSource = http.CanonicalHeaderKey("X-Ovh-Gateway-Source")
)

// LoggingMiddleware logs before and after incoming gin requests
func LoggingMiddleware(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if log == nil {
			c.Next()
			return
		}
		fields := logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"token":  c.Request.Header.Get(headerSource),
		}
		log.WithFields(fields).Debug("incoming request")

		c.Next()

		log.WithFields(fields).WithField("status", c.Writer.Status()).Info()

		for _, err := range c.Errors.Errors() {
			if err != RestErrorCreated.Error() {
				log.WithFields(fields).Error(err)
			}
		}
	}
}
