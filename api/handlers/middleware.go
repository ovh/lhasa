package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/security"
	"github.com/sirupsen/logrus"
)

// KeyRoles is the gin context key where roles are stored
const (
	KeyRoles = "roles"
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
			"source_token": c.GetHeader(security.HeaderGatewaySource),
			"request_id":   c.GetHeader(security.HeaderRequestID),
			"remote_user":  c.GetHeader(security.HeaderRemoteUser),
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

// HasOneRoleOf returns a gin handler that checks the request against the given role
func HasOneRoleOf(roles ...security.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesRaw, ok := c.Get(KeyRoles)
		rolePolicy, castok := rolesRaw.(security.RolePolicy)
		if !ok && !castok && rolePolicy != nil {
			// Bypass roles check if the configuration has not been properly set
			c.Next()
			return
		}

		if rolePolicy.HasOneRoleOf(roles...) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
	}
}

// AuthMiddleware returns a middleware that populate the Gin Context with security data
func AuthMiddleware(policy security.Policy) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(KeyRoles, security.BuildRolePolicy(policy, c.Request))
		c.Next()
	}
}
