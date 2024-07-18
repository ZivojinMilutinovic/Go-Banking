package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func SetupServer(nc *nats.Conn) {
	router := gin.Default()
	usersController := InitializeController()
	router.POST("create-user", func(ginContext *gin.Context) {
		usersController.CreateUser(ginContext)
	})

	router.GET("user-balance", func(ginContext *gin.Context) {
		usersController.Balance(ginContext, nc)
	})

	router.Run(":8080")
}
