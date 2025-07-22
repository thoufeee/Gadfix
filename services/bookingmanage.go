package services

import (
	"fmt"
	"gadfix/constansts"
	"gadfix/db"
	"gadfix/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// user servicing  book
func UserServiceRequest(c *gin.Context) {
	id := c.Param("id")
	serviceid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}

	var input struct {
		Address  string `json:"address" binding:"required"`
		Landmark string `json:"landmark" binding:"required"`
		Street   string `json:"street" binding:"required"`
		City     string `json:"city" binding:"required"`
		State    string `json:"state" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill blanks"})
		return
	}

	userid := c.MustGet("userid").(uint)

	var user models.User
	if err := db.DB.First(&user, &userid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "user not found"})
		return
	}

	address := models.UserAddress{
		UserID:   userid,
		Address:  input.Address,
		Landmark: input.Landmark,
		Street:   input.Street,
		City:     input.Street,
		State:    input.State,
	}

	if err := db.DB.Create(&address).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to store address"})
		return
	}

	// to get price of service
	var service models.Service

	if err := db.DB.First(&service, serviceid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "service not found"})
		return
	}

	price, _ := strconv.ParseFloat(service.Price, 64)

	booking := models.Booking{
		UserID:        userid,
		ServiceID:     uint(serviceid),
		Status:        constansts.StatusPending,
		PaymentStatus: "Pending",
		Amount:        price,
	}

	if err := db.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Proceed to Payment",
		"booking details": booking,
	})
}

// confirm payment
func ConfirmPayment(c *gin.Context) {
	id := c.Param("id")
	booking_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}

	var input struct {
		PaymentMode string `json:"mode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "choose a paymnet mode"})
		return
	}

	var booking models.Booking

	if err := db.DB.First(&booking, booking_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "booking not found"})
		return
	}

	paymentid := fmt.Sprintf("SIM-PAY-%d", time.Now().UnixNano())

	booking.PaymentId = paymentid
	booking.PaymentStatus = "Paid"
	booking.PaymentMode = input.PaymentMode

	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "payment failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Payment Successful",
		"msg":       "Booking Confirmed",
		"payemntid": paymentid,
	})
}

// admin approve service
func AdminApprove(c *gin.Context) {
	id := c.Param("id")
	booking_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}
	var booking models.Booking

	if err := db.DB.First(&booking, booking_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "booking not found"})
		return
	}
	booking.Status = constansts.StatusApproved
	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to approve booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "booking approved",
		"status":  booking,
	})
}

// admin assign staff
func AdminAssignStaff(c *gin.Context) {
	id := c.Param("id")
	booking_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}
	var data struct {
		StaffId uint `json:"staffid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "staff id required"})
		return
	}

	var booking models.Booking

	if err := db.DB.First(&booking, booking_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "booking not found"})
		return
	}

	booking.StaffID = &data.StaffId
	booking.Status = constansts.StatusAssigned

	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to assign staff"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "staff assigned",
	})
}

// staff Accept service
func StaffAcceptService(c *gin.Context) {
	id := c.Param("id")
	booking_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}
	staff_id := c.MustGet("userid").(uint)

	var booking models.Booking

	if err := db.DB.First(&booking, booking_id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "booking not found"})
		return
	}

	if booking.StaffID == nil || *booking.StaffID != staff_id {
		c.JSON(http.StatusForbidden, gin.H{"err": "not authorized to accept this booking"})
		return
	}

	booking.Status = constansts.StatusAccepted
	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to accept booking"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfuly accepted Booking"})
}

// Admin send response to user booking confirmed
func BookingConfirmed(c *gin.Context) {
	id := c.Param("id")
	booking_id, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}

	var data struct {
		PickupTime time.Time `json:"pickuptime" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill blanks"})
		return
	}

	var booking models.Booking

	if err := db.DB.First(&booking, booking_id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "booking not found"})
		return
	}

	booking.PickupTime = &data.PickupTime
	booking.Status = constansts.StatusPickupScheduled

	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to update booking"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfuly Send respones to User",
		"booking": booking,
	})
}

// booking cancel for user
func BookingCancel(c *gin.Context) {
	id := c.Param("id")
	bookingid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}

	var booking models.Booking
	if err := db.DB.First(&booking, bookingid).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "booking not found"})
		return
	}

	if err := db.DB.Unscoped().Delete(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to cancel booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfuly Canceled Booking"})

}

// staff picking product from user
func PickingConfirmed(c *gin.Context) {
	staff_id := c.MustGet("userid").(uint)
	id := c.Param("id")

	booking_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}

	var booking models.Booking
	if err := db.DB.First(&booking, booking_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "booking not found"})
		return
	}

	if booking.StaffID == nil || *booking.StaffID != staff_id {
		c.JSON(http.StatusBadRequest, gin.H{"err": "you are not assigned for this task"})
		return
	}

	time := time.Now()
	booking.PickedUpAt = &time
	booking.Status = constansts.StatusPicked

	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to confrim pickup"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery Accepetd"})
}

// Admin sets delivery time
func DeliveryTime(c *gin.Context) {
	id := c.Param("id")
	booking_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}

	var data struct {
		DeliveryTime time.Time `json:"deliverytime" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "fill blanks"})
		return
	}

	var booking models.Booking
	if err := db.DB.First(&booking, booking_id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "booking not found"})
		return
	}

	booking.DeliveryTime = &data.DeliveryTime
	booking.Status = constansts.StatusScheduled

	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "booking scheduled failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfuly scheduled time",
		"time":    booking.DeliveryTime,
	})
}

// delivery completed
func DeliveryCompleted(c *gin.Context) {
	id := c.Param("id")
	staff_id := c.MustGet("userid").(uint)

	booking_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}

	var booking models.Booking

	if err := db.DB.First(&booking, booking_id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "booking not found"})
		return
	}

	if booking.StaffID == nil || *booking.StaffID != staff_id {
		c.JSON(http.StatusForbidden, gin.H{"err": "you are not assigned for this booking"})
		return
	}

	time := time.Now()
	booking.DeliveredAt = &time
	booking.Status = constansts.StatusCompleted

	if err := db.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to save delivery"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Delivery Complete"})

}

// booking history
func BookingHistory(c *gin.Context) {
	var booking []models.Booking

	if err := db.DB.Find(&booking).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "failed to fetch booking history"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// sending staffs detalis to user
func BookingDetailsToUser(c *gin.Context) {
	id := c.Param("id")
	booking_id, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid id"})
		return
	}

	var booking models.Booking
	if err := db.DB.First(&booking, booking_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "booking not found"})
		return
	}

	var staff models.Staff
	if err := db.DB.First(&staff, booking.StaffID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "staff not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"staff_id":    staff.ID,
		"staff_name":  staff.FirstName,
		"status":      booking.Status,
		"pickup_time": booking.PickupTime,
	})

}

// sending user details to staff
func UserDetailsToStaff(c *gin.Context) {
	staff_id := c.MustGet("userid").(uint)

	var booking models.Booking
	if err := db.DB.Where("staff_id = ?", staff_id).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "no active booking found"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, booking.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "user not found"})
		return
	}

	var address models.UserAddress
	if err := db.DB.Where("user_id = ?", user.ID).First(&address).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"user_name":  user.FirstName + " " + user.SecondName,
			"user_phone": user.Phone,
			"address":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_name":  user.FirstName + " " + user.SecondName,
		"user_phone": user.Phone,
		"address":    address,
	})
}
