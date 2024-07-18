package controllers

import (
	"encoding/json"
	"net/http"
	"time"
	"users/api"
	"users/models"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

type UsersController struct {
	userApi api.UserInterface
}

func InitializeController() UsersController {
	return UsersController{
		userApi: api.GetUserApi(),
	}
}

func (controller UsersController) CreateUser(ginContext *gin.Context) {
	var userPayload models.UserPayload

	if err := ginContext.ShouldBindJSON(&userPayload); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, apiError := controller.userApi.CreateUser(&userPayload)

	if apiError != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": apiError.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "User created successfully!", "user": user})
}

func (controller UsersController) Balance(ginContext *gin.Context, nc *nats.Conn) {

	email := ginContext.Query("email")
	if email == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Email query parameter is required"})
		return
	}

	msg, err := nc.Request("get_balance", []byte(email), 2*time.Second)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch balance:" + err.Error()})
		return
	}

	var userResponse models.UserResponse
	if err := json.Unmarshal(msg.Data, &userResponse); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse balance response:" + err.Error()})
		return
	}

	if userResponse.Error != "" {
		ginContext.JSON(http.StatusNotFound, gin.H{"error": userResponse.Error})
		return
	}

	ginContext.JSON(http.StatusOK, userResponse)
}
