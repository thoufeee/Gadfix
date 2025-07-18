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

// staff signup
func StaffSignup(c *gin.Context) {
	var data struct {
		FirstName    string `json:"firstname"`
		SecondName   string `json:"secondname" binding:"required"`
		Email        string `json:"email" binding:"required"`
		Password     string `json:"password" binding:"required,min=6"`
		Phone        string `json:"phone" binding:"required,len=10"`
		IdentityCard string `json:"cardnumber" binding:"required,len=12"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	//  email check
	var existing models.Staff
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
	data.IdentityCard = strings.TrimSpace(data.IdentityCard)

	newstaff := models.Staff{
		FirstName:    data.FirstName,
		SecondName:   data.SecondName,
		Email:        data.Email,
		Password:     hashedpass,
		Phone:        data.Phone,
		IdentityCard: data.IdentityCard,
		Role:         constansts.RoleStaff,
	}

	if err := db.DB.Create(&newstaff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "could not create staff account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "successfuly registered"})
}
