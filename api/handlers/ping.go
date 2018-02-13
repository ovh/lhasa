package handlers

import "github.com/gin-gonic/gin"

// PingHandler provides monitoring status on a rest endpoint
func PingHandler(c *gin.Context) (string, error) {
	return "OK", nil
}
