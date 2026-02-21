package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Login        string     `gorm:"unique;not null"`
	HashPassword string     `gorm:"not null"`
	BaseCurrence string     `gorm:"not null"`
	Tokens       []Token    `gorm:"foreignKey:UserID"`
	Accounts     []*Account `gorm:"many2many:user_accounts;"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
