package controllers

import (
	"gadfix/constansts"
	"gadfix/db"
	"gadfix/models"
	"gadfix/utlis"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// signup for user

func Signup(c *gin.Context) {

	var data struct {
		FirstName  string `json:"firstname" binding:"required"`
		SecondName string `json:"secondname" binding:"required"`
		Email      string `json:"email" binding:"required"`
		Password   string `json:"password" binding:"required,min=6"`
		Phone      string `json:"phone" binding:"required,len=10"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
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
	data.Phone = strings.TrimSpace(data.Phone)

	newuser := models.User{
		FirstName:  data.FirstName,
		SecondName: data.SecondName,
		Email:      data.Email,
		Password:   hashedpass,
		Phone:      data.Phone,
		Role:       constansts.RoleUser,
	}

	if err := db.DB.Create(&newuser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "could not create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "successfuly registered"})
}
