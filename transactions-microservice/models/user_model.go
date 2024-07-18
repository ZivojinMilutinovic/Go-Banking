package models

import "time"

type User struct {
	UserId    int       `gorm:"primaryKey;column:user_id;autoIncrement"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null"`
	Balance   float64   `gorm:"type:float;not null;default:0"`
}

type UserBalanceResponse struct {
	Email   string  `json:"email"`
	Balance float64 `json:"balance"`
	Error   string  `json:"error,omitempty"`
}
