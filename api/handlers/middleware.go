package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ovh/lhasa/api/hateoas"
	"github.com/ovh/lhasa/api/security"
	"github.com/sirupsen/logrus"
)

// Gin context keys
const (
	KeyRoles    = "roles"
	KeyLogEntry = "logger"
)

// LoggingMiddleware logs before and after incoming gin requests
func LoggingMiddleware(logHeaders []string, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if log == nil {
			c.Next()
			return
		}
		entry := log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		})
		for _, h := range logHeaders {
			header := http.CanonicalHeaderKey(h)
			if header != "" {
				entry = entry.WithField(header, c.GetHeader(header))
			}
		}
		c.Set(KeyLogEntry, entry)

		entry.Debug("incoming request")
		startTime := time.Now()

		c.Next()

		entry.WithFields(logrus.Fields{
			"status":   c.Writer.Status(),
			"duration": time.Since(startTime).Seconds(),
		}).Debug("done")

		for _, err := range c.Errors.Errors() {
			if err != hateoas.ErrorCreated.Error() {
				entry.Error(err)
			}
		}
	}
}

// HasOne returns a gin handler that checks the request against the given role
func HasOne(roles ...security.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		rolesRaw, ok := c.Get(KeyRoles)
		rolePolicy, castok := rolesRaw.(security.RolePolicy)
		log := GetLogger(c)

		if !ok || !castok {
			if log != nil {
				log.Warnf("unusable security policy, please check your configuration")
				abortUnauthorized(c)
				return
			}
		}

		if rolePolicy.HasOne(roles...) {
			c.Next()
			return
		}
		if log != nil {
			log.WithField("expectedRoles", roles).WithField("roles", rolePolicy.ToSlice()).Info("unauthorized access")
		}
		abortUnauthorized(c)
	}
}

func abortUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access"})
}

// AuthMiddleware returns a middleware that populate the Gin Context with security data
func AuthMiddleware(policy security.CompiledPolicy) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := security.BuildRolePolicy(policy, c.Request)
		c.Set(KeyRoles, roles)
		log := GetLogger(c)
		if log != nil {
			log.WithField("roles", roles.ToSlice()).Debugf("user has roles")
		}
		c.Next()
	}
}

// GetLogger returns the request logger from the gin context, nil if it does not exist
func GetLogger(c *gin.Context) *logrus.Entry {
	logRaw, found := c.Get(KeyLogEntry)
	log, castok := logRaw.(*logrus.Entry)
	if !found || !castok {
		return nil
	}
	return log
}
