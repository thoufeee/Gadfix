package controllers

import (
	"gadfix/db"
	"gadfix/models"
	"gadfix/utlis"
	"net/http"

	"github.com/gin-gonic/gin"
)

// login
func Login(c *gin.Context) {
	var data struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill all fileds"})
		return
	}

	// checking email for user
	var user models.User
	if err := db.DB.Where("email = ?", data.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid email or pass"})
		return
	}

	// block check
	if user.Block && user.Role != "0" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "your account is blocked"})
		return
	}

	// checking password
	if !utlis.CheckHash(user.Password, data.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid email or password"})
		return
	}

	// generate token
	token, err := utlis.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to create accesstoken"})
		return
	}

	// refresh token
	refresh, err := utlis.Refresh(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to create access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"res":     "Successfuly loged",
		"access":  token,
		"refresh": refresh,
		"role":    user.Role,
	})

}

// // forgott password
// func ForgetPassword(password string) {
//          id :=
// }
