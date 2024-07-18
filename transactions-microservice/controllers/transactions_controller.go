package controllers

import (
	"net/http"
	"strconv"
	"time"
	"transactions/api"
	"transactions/conn"
	"transactions/models"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

type TransactionsController struct {
	userApi api.UserInterface
}

func InitializeController(nc *nats.Conn) TransactionsController {
	return TransactionsController{
		userApi: api.GetUserApi(nc),
	}
}

func (transactionsController *TransactionsController) AddFunds(ginContext *gin.Context) {
	userId, err := strconv.Atoi(ginContext.Param("user_id"))
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := ginContext.ShouldBindJSON(&req); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tx := conn.GetDB().Begin()

	if tx.Error != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	var user models.User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, "user_id = ?", userId).Error; err != nil {
		tx.Rollback()
		ginContext.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Balance += req.Amount

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	transaction := models.Transaction{
		UserId:    userId,
		Amount:    req.Amount,
		Type:      "credit",
		CreatedAt: time.Now().UTC(),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"updated_balance": user.Balance})
}

func (transactionsController *TransactionsController) TransferFunds(ginContext *gin.Context) {
	var req struct {
		FromUserId       int     `json:"from_user_id" binding:"required"`
		ToUserId         int     `json:"to_user_id" binding:"required"`
		AmountToTransfer float64 `json:"amount_to_transfer" binding:"required"`
	}

	if err := ginContext.ShouldBindJSON(&req); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.AmountToTransfer <= 0 {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Transfer amount must be positive"})
		return
	}

	tx := conn.GetDB().Begin()

	if tx.Error != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	var fromUser, toUser models.User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&fromUser, "user_id = ?", req.FromUserId).Error; err != nil {
		tx.Rollback()
		ginContext.JSON(http.StatusNotFound, gin.H{"error": "Sender not found"})
		return
	}
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&toUser, "user_id = ?", req.ToUserId).Error; err != nil {
		tx.Rollback()
		ginContext.JSON(http.StatusNotFound, gin.H{"error": "Receiver not found"})
		return
	}

	if fromUser.Balance < req.AmountToTransfer {
		tx.Rollback()
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	fromUser.Balance -= req.AmountToTransfer
	toUser.Balance += req.AmountToTransfer

	if err := tx.Save(&fromUser).Error; err != nil {
		tx.Rollback()
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sender balance"})
		return
	}

	if err := tx.Save(&toUser).Error; err != nil {
		tx.Rollback()
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receiver balance"})
		return
	}

	debitTransaction := models.Transaction{
		UserId:    fromUser.UserId,
		Amount:    -req.AmountToTransfer,
		Type:      "debit",
		CreatedAt: time.Now(),
	}
	creditTransaction := models.Transaction{
		UserId:    toUser.UserId,
		Amount:    req.AmountToTransfer,
		Type:      "credit",
		CreatedAt: time.Now(),
	}

	if err := tx.Create(&debitTransaction).Error; err != nil {
		tx.Rollback()
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record debit transaction"})
		return
	}

	if err := tx.Create(&creditTransaction).Error; err != nil {
		tx.Rollback()
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record credit transaction"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"from_user_balance": fromUser.Balance,
		"to_user_balance":   toUser.Balance,
	})
}
