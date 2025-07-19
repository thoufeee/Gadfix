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

	// phone number check
	if err := db.DB.Where("phone =?", existing.Phone).First(&existing).Error; err == nil {
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

// staff profile update
func StaffProfileUpdate(c *gin.Context) {
	staff_id := c.MustGet("userid").(uint)

	var staff models.Staff
	if err := db.DB.First(&staff, staff_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "staff not found"})
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
		staff.Phone = *input.Phone
	}

	if input.FirstName != nil {
		staff.FirstName = *input.FirstName
	}

	if input.SecondName != nil {
		staff.SecondName = *input.SecondName
	}

	if err := db.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfuly updated profile"})
}

// staff forgott password
func StaffForgotPassword(c *gin.Context) {
	var input struct {
		Phone    string `json:"phone" binding:"required,len=10"`
		Password string `josn:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	var staff models.Staff
	if err := db.DB.Where("phone = ?", input.Phone).First(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "account with this phone number not found"})
		return
	}

	hashpass, err := utlis.Hash(input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "password not hashed"})
		return
	}

	staff.Password = hashpass

	if err := db.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password successfuly updated"})
}
