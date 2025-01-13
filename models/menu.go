package models

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	Name        string `gorm:"notnull"`
	Description string
	ImageFileID *string
}
