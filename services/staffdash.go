package services

import (
	"gadfix/db"
	"gadfix/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// staff dash
func StaffDash(c *gin.Context) {
	staffID := c.MustGet("userid").(uint)

	var staff models.Staff
	if err := db.DB.First(&staff, staffID).Error; err != nil {

		if strings.Contains(c.GetHeader("Accept"), "application/json") {

			c.JSON(http.StatusBadRequest, gin.H{"err": "staff account not found"})
		} else {

			c.HTML(http.StatusBadRequest, "error.html", gin.H{"err": "staff not found"})
		}
		return
	}

	if strings.Contains(c.GetHeader("Accept"), "application/json") {
		c.JSON(http.StatusOK, gin.H{"profile": staff})
	} else {
		c.HTML(http.StatusOK, "staff.html", gin.H{
			"title":   "staff",
			"account": staff,
		})

	}
}
