package services

import (
	"gadfix/db"
	"gadfix/models"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		Name         *string `json:"name"`
		Price        *string `json:"price"`
		Description  *string `json:"description"`
		Category     *string `json:"category"`
		Duration     *string `json:"duration"`
		ServiceImage *string `json:"url"`
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

	if input.Name != nil {
		service.Name = *input.Name
	}
	if input.Price != nil {
		service.Price = *input.Price
	}
	if input.Description != nil {
		service.Description = *input.Description
	}
	if input.Category != nil {
		service.Category = *input.Category
	}
	if input.Duration != nil {
		service.Duration = *input.Duration
	}
	if input.ServiceImage != nil {
		service.ServiceImage = *input.ServiceImage
	}

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
		c.JSON(http.StatusNotFound, gin.H{"err": "service not found"})
		return
	}

	// checking service has active booking
	var count int64
	db.DB.Model(&models.Booking{}).Where("service_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "service has active bookings and cannot be deleted"})
		return
	}

	if err := db.DB.Unscoped().Delete(&service).Error; err != nil {
		log.Println("delete error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"err": "service deletion failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"res": "successfuly deleted service"})
}

// search services
func SearchService(c *gin.Context) {
	key := c.Query("search")
	key = strings.TrimSpace(key)

	var services []models.Service

	if key == "" {
		c.JSON(http.StatusOK, []models.Service{})
		return
	}

	pattern := "%" + key + "%"

	if err := db.DB.Where("LOWER(name) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", pattern, pattern).
		Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to search services"})
		return
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
