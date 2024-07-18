package controllers

import (
	"github.com/gin-gonic/gin"
)

func SetupServer() {
	router := gin.Default()
	transactionController := InitializeController(nil)
	router.PATCH("/users/:user_id/balance", func(ginContext *gin.Context) {
		transactionController.AddFunds(ginContext)
	})
	router.PATCH("/transfer", func(ginContext *gin.Context) {
		transactionController.TransferFunds(ginContext)
	})

	router.Run(":8081")
}
