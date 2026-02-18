package models

import (
	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	UserID       uint   `gorm:"not null;index"`
	RefreshToken string `gorm:"not null"`
}

func (Token) TableName() string {
	return "tokens"
}
