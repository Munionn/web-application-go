package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Users []*User `gorm:"many2many:user_accounts;"`
	Name  string  `gorm:"not null"`
}

func (Account) TableName() string {
	return "accounts"
}
