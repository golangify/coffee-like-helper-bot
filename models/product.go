package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	MenuID      uint
	Menu        *Menu  `gorm:"foreignKey:MenuID"`
	Name        string `gorm:"notnull"`
	Description string
	ImageFileID *string
}
