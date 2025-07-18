package models

import (
	"time"

	"gorm.io/gorm"
)

// user model
type User struct {
	gorm.Model
	FirstName  string        `json:"firstname"`
	SecondName string        `json:"secondname"`
	Email      string        `json:"email" gorm:"unique"`
	Password   string        `json:"password"`
	Role       string        `json:"role"`
	Phone      string        `json:"phone"`
	Block      bool          `gorm:"default:false"`
	Address    []UserAddress `json:"address" gorm:"foreignKey:UserID"`
	Bookings   []Booking     `json:"bookings" gorm:"foreignKey:UserID"`
}

// staff model
type Staff struct {
	gorm.Model
	FirstName    string `json:"firstname"`
	SecondName   string `json:"secondname" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
	IdentityCard string `json:"cardnumber" binding:"required,len=12"`
	Phone        string `json:"phone" binding:"required,len=10"`
	Role         string `json:"role"`
	Block        bool   `gorm:"default:false"`
}

// address model
type UserAddress struct {
	gorm.Model
	UserID   uint   `json:"userid"`
	Address  string `json:"address"`
	Landmark string `json:"landmark"`
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
}

// bookings
type Booking struct {
	gorm.Model
	UserID    uint   `json:"userid"`
	ServiceID uint   `json:"serviceid"`
	StaffID   *uint  `json:"staffid"`
	Status    string `json:"status"`

	PaymentStatus string  `json:"paymentstatus"`
	PaymentId     string  `json:"payment_id"`
	Amount        float64 `json:"amount"`
	PaymentMode   string  `json:"paymentmode"`

	PickupTime *time.Time // setting admin
	PickedUpAt *time.Time // setting staff

	DeliveryTime *time.Time //  setting admin
	DeliveredAt  *time.Time // settig staff

	CreatedAt *time.Time // by user
	UpdatedAt *time.Time //  by user
}

// services
type Service struct {
	gorm.Model
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Price        string    `json:"price"`
	Description  string    `json:"description"`
	Duration     string    `json:"duration"`
	ServiceImage string    `json:"url" binding:"required"`
	Booking      []Booking `json:"booking" gorm:"foreignKey:ServiceID"`
}
