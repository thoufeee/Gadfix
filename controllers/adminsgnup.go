package controllers

import (
	"gadfix/constansts"
	"gadfix/db"
	"gadfix/models"
	"gadfix/utlis"
	"net/http"

	"github.com/gin-gonic/gin"
)

// admin signup
func AdminSignup(c *gin.Context) {
	var data struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill all fileds"})
		return
	}

	//  email check
	var existing models.User
	if err := db.DB.Where("email = ?", data.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "email already registered"})
		return
	}

	// hashing
	hashedpass, err := utlis.Hash(data.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "password not hashed"})
		return
	}

	newadmin := models.User{
		Email:    data.Email,
		Password: hashedpass,
		Role:     constansts.RoleAdmin,
	}

	if err := db.DB.Create(&newadmin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "could not create admin account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "successfuly registered"})
}
