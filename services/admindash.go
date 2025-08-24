package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// admin dashboard
func AdminDash(c *gin.Context) {
	email, _ := c.Get("email")
	c.JSON(http.StatusOK, gin.H{"res": email})
}
