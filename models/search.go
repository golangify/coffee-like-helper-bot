package models

import (
	"gorm.io/gorm"
)

type Search struct {
	gorm.Model
	Text   string
	UserID uint
}
