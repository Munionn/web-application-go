package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Users        []*User `gorm:"many2many:user_accounts;"`
	Name         string  `gorm:"not null"`
	BaseCurrence string  `gorm:"not null"`
	Balance      float64 `json:"balance" gorm:"type:decimal;default:0"`
}

func (Account) TableName() string {
	return "accounts"
}
