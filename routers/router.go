package routers

import (
	"gadfix/constansts"
	"gadfix/controllers"

	"gadfix/middleware"
	"gadfix/services"

	"github.com/gin-gonic/gin"
)

// public
func Public(c *gin.Engine) {
	api := c.Group("/api")
	api.POST("/signup", controllers.Signup)
	api.POST("/login", controllers.Login)
	api.POST("/logout", controllers.Logout)
	api.PUT("/forgottpassworduser", controllers.UserForgotPassword)

	api.POST("/staffsignup", controllers.StaffSignup)
	api.POST("/stafflogin", controllers.StaffLogin)
	api.PUT("/forgottpasswordstaff", controllers.StaffForgotPassword)

	api.POST("/adminsignup", controllers.AdminSignup)

	api.GET("/get", services.ServiceListing)
	api.GET("/search", services.SearchService)
}

// user only route
func UserRoute(c *gin.Engine) {
	user := c.Group("/user")
	user.Use(middleware.Auth("1"))
	user.GET("/dash", services.UserDash)

	user.POST("/booking/:id", services.UserServiceRequest)
	user.POST("/payment/:id", services.ConfirmPayment)
	user.DELETE("/bookingcancel/:id", services.BookingCancel)
	user.GET("/staffdetails/:id", services.BookingDetailsToUser)
	user.PUT("/profileupdate", controllers.UpdateUserProfile)
	user.GET("/history", services.UserBookingHistory)
}

// admin only route
func AdminRoute(c *gin.Engine) {
	admin := c.Group("/admin")
	admin.Use(middleware.Auth(constansts.RoleAdmin))
	admin.GET("/dash", services.AdminDash)
	admin.POST("/adminlogout", controllers.AdminSignout)

	admin.GET("/users", services.UserDetails)
	admin.GET("/users/length", services.UsersTotalLength)
	admin.PUT("/block/:id", services.BlockUsers)
	admin.PUT("/unblock/:id", services.UnblockUSers)
	admin.POST("/createuser", services.CreateUsers)
	admin.DELETE("/delete/:id", services.DeleteUsers)
	admin.PUT("/updateuser/:id", services.UpdateUsers)

	admin.GET("/service/length", services.ServiceLength)
	admin.POST("/service/add", services.CreateService)
	admin.PUT("/service/update/:id", services.UpdateServices)
	admin.DELETE("/service/delete/:id", services.DeleteServices)

	admin.GET("/staff/length", services.StaffTotalLength)
	admin.GET("/staffs", services.ListStaff)
	admin.POST("/createstaff", services.CreateStaff)
	admin.PUT("/staffblock/:id", services.BlockStaff)
	admin.PUT("/staffunblock/:id", services.UnBlockStaff)
	admin.DELETE("/staffdelete/:id", services.DeleteStaff)
	admin.PUT("/updatestaff/:id", services.UpdateStaffProfile)

	admin.PUT("/bookingapproved/:id", services.AdminApprove)
	admin.POST("/assignstaff/:id", services.AdminAssignStaff)
	admin.POST("/bookingconfirmed/:id", services.BookingConfirmed)
	admin.POST("/deliverytime/:id", services.DeliveryTime)
	admin.GET("/bookinghistory", services.BookingHistory)

}

// staff only route
func StaffRoute(c *gin.Engine) {
	staff := c.Group("/staff")
	staff.Use(middleware.Auth(constansts.RoleStaff))
	staff.GET("/dash", services.StaffDash)
	staff.POST("/stafflogout", controllers.StaffLogout)

	staff.PUT("/acceptservice/:id", services.StaffAcceptService)
	staff.PUT("/deliveryconfirmed/:id", services.PickingConfirmed)
	staff.PUT("/deliverycompleted/:id", services.DeliveryCompleted)
	staff.GET("/useraddress", services.UserDetailsToStaff)
	staff.PUT("/profileupdate", controllers.StaffProfileUpdate)
}
