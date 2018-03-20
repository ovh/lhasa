package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// PingHandler provides monitoring status on a rest endpoint
func PingHandler(db *sql.DB) func(c *gin.Context) (string, error) {
	return func(c *gin.Context) (string, error) {
		err := db.Ping()
		if err != nil {
			return "", err
		}
		return "OK", nil
	}
}
