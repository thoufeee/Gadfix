package services

import (
	"gadfix/constansts"
	"gadfix/db"
	"gadfix/models"
	"gadfix/utlis"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// all user details
func UserDetails(c *gin.Context) {
	var data []models.User

	if err := db.DB.Where("role = ?", 1).Find(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// block users
func BlockUsers(c *gin.Context) {
	userid := c.Param("id")
	var user models.User

	if err := db.DB.First(&user, userid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}
	user.Block = true
	db.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"res": "successfuly blocked"})
}

// unblock users
func UnblockUSers(c *gin.Context) {
	userid := c.Param("id")

	var user models.User

	if err := db.DB.First(&user, userid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}
	user.Block = false
	db.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"res": "successfuly unblocked"})
}

// delete users
func DeleteUsers(c *gin.Context) {
	userid := c.Param("id")

	id, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}
	var user models.User

	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}

	if err := db.DB.Unscoped().Delete(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": "successfuly deleted"})
}

// length of total users
func UsersTotalLength(c *gin.Context) {
	var users []models.User

	if err := db.DB.Where("role = ?", 1).First(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to fetch userlength"})
		return
	}
	c.JSON(http.StatusOK, len(users))
}

// creating users
func CreateUsers(c *gin.Context) {
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

	var user models.User
	// email checking
	if err := db.DB.Where("email = ?", data.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "email already existed"})
		return
	}

	hashedpass, err := utlis.Hash(data.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to hash pass"})
		return
	}

	phone := strings.TrimSpace(data.Phone)

	newuser := models.User{
		FirstName:  data.FirstName,
		SecondName: data.SecondName,
		Email:      data.Email,
		Password:   hashedpass,
		Phone:      phone,
		Role:       constansts.RoleUser,
	}

	if err := db.DB.Create(&newuser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to create user account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mess": "successfuly created user",
		"data": newuser,
	})

}
