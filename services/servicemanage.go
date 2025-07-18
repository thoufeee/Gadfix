package services

import (
	"gadfix/db"
	"gadfix/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// listing services
func ServiceListing(c *gin.Context) {
	var data []models.Service
	if err := db.DB.Find(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to get services"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// create services
func CreateService(c *gin.Context) {
	var data models.Service
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	if err := db.DB.Create(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to create services"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "service created successfuly"})
}

// update services
func UpdateServices(c *gin.Context) {
	userid := c.Param("id")
	id, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid service id"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Price       string `json:"price"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid input"})
		return
	}

	var service models.Service
	if err := db.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "service not found"})
		return
	}

	service.Name = input.Name
	service.Price = input.Price
	service.Description = input.Description

	if err := db.DB.Save(&service).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to change services"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": "successfuly updated service"})
}

// delete services
func DeleteServices(c *gin.Context) {
	userid := c.Param("id")
	id, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id input"})
		return
	}

	var service models.Service
	if err := db.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "service not found"})
		return
	}

	if err := db.DB.Unscoped().Delete(&service).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "service deletion failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": "successfuly deleted service"})
}

// search services
func SearchService(c *gin.Context) {
	key := c.Query("search")
	var services []models.Service

	if key != "" {
		pattern := "%" + key + "%"
		if err := db.DB.Where("name LIKE ? OR description LIKE ?", pattern, pattern).Find(&services).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "failed to fetch services"})
			return
		}
	} else {
		if err := db.DB.Find(&services).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "service not found"})
			return
		}
	}
	c.JSON(http.StatusOK, services)
}

// length of services
func ServiceLength(c *gin.Context) {
	var services []models.Service
	if err := db.DB.Find(&services).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to get services"})
		return
	}

	c.JSON(http.StatusOK, len(services))
}
