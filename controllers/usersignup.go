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

	// phone number check
	if err := db.DB.Where("phone = ?", existing.Phone).First(&existing); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "phone number already registered"})
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

// update user profile
func UpdateUserProfile(c *gin.Context) {

	user_id := c.MustGet("userid").(uint)

	var user models.User
	if err := db.DB.First(&user, user_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}

	var input struct {
		FirstName  *string `json:"firstname"`
		SecondName *string `json:"secondname"`
		Phone      *string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	if input.Phone != nil {
		if len(*input.Phone) != 10 {
			c.JSON(http.StatusBadRequest, gin.H{"err": "password must be 10 digits"})
			return
		}
		user.Phone = *input.Phone
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}

	if input.SecondName != nil {
		user.SecondName = *input.SecondName
	}

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfuly"})
}

// forgott password for user
func UserForgotPassword(c *gin.Context) {
	var input struct {
		Phone    string `json:"phone" binding:"required,len=10"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	var user models.User
	if err := db.DB.Where("phone = ?", input.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "account with this phone number not found"})
		return
	}

	hashpass, err := utlis.Hash(input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "password not hashed"})
		return
	}

	if err := db.DB.Model(&user).Update("password", hashpass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfuly updated password"})
}
