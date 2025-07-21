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

	var booking models.Booking

	if err := db.DB.First(&booking, staff_id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "booking not found"})
		return
	}

	var user models.User

	if err := db.DB.First(&user, booking.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "user not found"})
		return
	}

	var address models.UserAddress

	if err := db.DB.First(&address, user.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "user address not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"staff":        staff,
		"booking_id":   booking.ID,
		"user_id":      user.ID,
		"user_name":    user.FirstName + " " + user.SecondName,
		"user_phone":   user.Phone,
		"user_address": address,
	})
}
