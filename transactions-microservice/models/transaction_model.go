package models

import "time"

type Transaction struct {
	ID        int       `gorm:"primaryKey"`
	UserId    int       `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	Type      string    `gorm:"type:varchar(10);not null"` // "credit" or "debit"
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
