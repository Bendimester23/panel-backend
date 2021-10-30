package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ultimatepanel2000/backend/db"
	"github.com/ultimatepanel2000/backend/router"

	"github.com/gin-contrib/cors"
)

func main() {
	log.Println("Starting")

	db.Connect()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(200, c.ClientIP())
	})

	router.InitAuth(r.Group("/auth"))
	router.InitUser(r.Group("/user"))

	defer db.Disconnect()

	r.Run(":8080")
}
