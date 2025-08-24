package services

import (
	"gadfix/constansts"
	"gadfix/db"
	"gadfix/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// user dashboard
func UserDash(c *gin.Context) {

	user_id := c.MustGet("userid").(uint)

	//  find user
	var user models.User
	if err := db.DB.Find(&user, user_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}

	//  find booking
	var booking models.Booking
	if err := db.DB.Where("user_id = ? AND status != ?", user_id, constansts.StatusCompleted).Order("created_at DESC").
		First(&booking).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"account":    user,
			"booking_id": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account":    user,
		"booking_id": booking.ID,
	})
}
