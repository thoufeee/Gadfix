package controllers

import (
	"gadfix/db"
	"gadfix/models"
	"gadfix/utlis"
	"net/http"

	"github.com/gin-gonic/gin"
)

// staff login
func StaffLogin(c *gin.Context) {
	var data struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill all fileds"})
		return
	}

	// checking email
	var staff models.Staff

	if err := db.DB.Where("email = ?", data.Email).First(&staff).Error; err == nil {

		// check block
		if staff.Block {
			c.JSON(http.StatusBadRequest, gin.H{"err": "your account is blocked"})
			return
		}

		// checking password
		if !utlis.CheckHash(staff.Password, data.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"err": "invalid email or password"})
			return
		}

		// generate token
		token, err := utlis.Generate(staff.ID, staff.Email, staff.Role)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "failed to generate token"})
			return
		}

		// refresh token
		refresh, err := utlis.Refresh(staff.ID, staff.Email, staff.Role)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "failed to generate refresh token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"res": "successfuly loged in",
			"access token":  token,
			"refresh token": refresh,
			"role":          staff.Role,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"err": "invalid email or password"})
}
