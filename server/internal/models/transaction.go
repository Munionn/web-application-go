package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Amount           uint   `json:"amount" gorm:"type:integer;not null"`
	BaseCurrency     string `json:"base_currency" gorm:"type:varchar(10);not null"`
	Type             string `json:"type" gorm:"type:varchar(20);not null"`
	ShortDescription string `json:"short_description" gorm:"type:text"`

	UserID    uint `gorm:"not null;index" json:"user_id"`
	User      User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AccountID uint `gorm:"not null;index" json:"account_id"`
	Account   Account `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

func (Transaction) TableName() string {
	return "transactions"
}
