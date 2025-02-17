package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramID      int64 `gorm:"unique;notnull"`
	IsBarista       bool
	IsAdministrator bool
	IsBanned        bool
	FirstName       string
	LastName        string
	UserName        string
}
