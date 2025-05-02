package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Uuid     string `gorm:"size:100;uniqueIndex;not null"`
	Name     string `gorm:"size:100;not null"`
	Email    string `gorm:"size:255;uniqueIndex;not null"`
	Password string `gorm:"size:100;not null"`
}
