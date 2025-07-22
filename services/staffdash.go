package services

import (
	"gadfix/db"
	"gadfix/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// staff dash
func StaffDash(c *gin.Context) {
	staff_id := c.MustGet("userid").(uint)

	var staff models.Staff

	if err := db.DB.First(&staff, staff_id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "staff not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"staff": staff,
	})
}
