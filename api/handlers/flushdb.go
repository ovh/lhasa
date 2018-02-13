package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ovh/lhasa/api/models"
)

// FlushDBHandler empties the Database for testing purposes
func FlushDBHandler(c *gin.Context) (string, error) {
	if err := models.FlushApplications(); err != nil {
		return "NOK", err
	}
	return "OK", nil
}
