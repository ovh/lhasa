package handlers

import (
	"github.com/gin-gonic/gin"
)

// VersionHandler displays the current version number
func VersionHandler(version string) func(*gin.Context) (string, error) {
	return func(_ *gin.Context) (string, error) {
		return version, nil
	}
}
