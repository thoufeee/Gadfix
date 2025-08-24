package services

import (
	"gadfix/db"
	"gadfix/models"
	"gadfix/utlis"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// list all users
func ListStaff(c *gin.Context) {
	var staffs []models.Staff

	if err := db.DB.Find(&staffs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to get datas of staffs"})
		return
	}

	c.JSON(http.StatusOK, staffs)
}

// block staff
func BlockStaff(c *gin.Context) {
	userid := c.Param("id")
	id, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}
	var staff models.Staff
	if err := db.DB.Find(&staff, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "staff not found"})
		return
	}
	staff.Block = true
	if err := db.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to block"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": "successfuly blocked"})
}

// unblock staff
func UnBlockStaff(c *gin.Context) {
	userid := c.Param("id")
	id, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}
	var staff models.Staff
	if err := db.DB.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "staff not found"})
		return
	}

	staff.Block = false
	if err := db.DB.Save(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to block staff"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": "successfuly unblocked"})
}

// delete staff
func DeleteStaff(c *gin.Context) {
	userid := c.Param("id")
	id, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}
	var staff models.Staff

	if err := db.DB.First(&staff, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "staff not found"})
		return
	}

	if err := db.DB.Unscoped().Delete(&staff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "faiiled to delete staff"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "successfuly deleted staff"})
}

// total length of staffs
func StaffTotalLength(c *gin.Context) {
	var staffs []models.Staff

	if err := db.DB.First(&staffs).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to get length off staffs"})
		return
	}

	c.JSON(http.StatusOK, len(staffs))
}

// create staffs
func CreateStaff(c *gin.Context) {
	var data struct {
		FirstName    string `json:"firstname" binding:"required"`
		SecondName   string `json:"secondname" binding:"required"`
		Email        string `json:"email" binding:"required"`
		Password     string `json:"password" binding:"required,min=6"`
		Phone        string `json:"phone" binding:"required,len=10"`
		IdentityCard string `json:"cardnumber" binding:"required,len=12"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill all fileds"})
		return
	}

	var staff models.Staff
	// email check
	if err := db.DB.Where("email = ?", staff.Email).First(&staff).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "email already registered"})
		return
	}

	hashedpass, err := utlis.Hash(staff.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to hash password"})
		return
	}

	phone := strings.TrimSpace(data.Password)
	identitycard := strings.TrimSpace(data.IdentityCard)

	newstaff := models.Staff{
		FirstName:    data.FirstName,
		SecondName:   data.SecondName,
		Email:        data.Email,
		Password:     hashedpass,
		Phone:        phone,
		IdentityCard: identitycard,
	}

	if err := db.DB.Create(&newstaff).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to create staff account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mess":  "successfuly created staff",
		"staff": newstaff,
	})
}

// update staff profile
func UpdateStaffProfile(c *gin.Context) {
	id := c.Param("id")
	staff_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}

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
			c.JSON(http.StatusBadRequest, gin.H{"err": "phone number must be 10 digits"})
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
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to update staff details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfuly updated staff profile"})
}
