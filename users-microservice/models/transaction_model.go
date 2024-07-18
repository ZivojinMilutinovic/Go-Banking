package models

import "time"

type Transaction struct {
	ID        int       `gorm:"primaryKey"`
	UserId    int       `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
