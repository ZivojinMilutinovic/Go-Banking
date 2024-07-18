package models

import "time"

type User struct {
	UserId    int       `gorm:"primaryKey;column:user_id;autoIncrement"`
	Email     string    `gorm:"type:varchar(100);unique;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null"`
}
