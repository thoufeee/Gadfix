package services

import (
	"gadfix/db"
	"gadfix/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// user dashboard
func UserDash(c *gin.Context) {

	user_id := c.MustGet("userid").(uint)

	var user models.User
	if err := db.DB.Find(&user, user_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account": user,
	})
}
