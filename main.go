package main

import (
	"gadfix/db"
	"gadfix/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// database
	db.Connection()
	// redis
	db.ConnectRedis()

	r := gin.Default()

	r.Static("/static", "./templates")
	// r.LoadHTMLGlob("templates/*.html")

	routers.Public(r)
	routers.UserRoute(r)
	routers.StaffRoute(r)
	routers.AdminRoute(r)

	r.Run(":8080")

}
