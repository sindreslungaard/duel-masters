package db

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	ID              uint
	ReferenceString string
}
