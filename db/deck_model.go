package db

import "gorm.io/gorm"

type Deck struct {
	gorm.Model
	ID       uint     `json:"uid"`
	UserID   uint     `json:"owner"`
	Name     string   `json:"name"`
	Public   bool     `json:"public"`
	Standard bool     `json:"standard"`
	Cards    []string `json:"cards"`
}

type DeckCards struct {
	gorm.Model
	DeckID uint `gorm:"primaryKey"`
	CardID uint `gorm:"primaryKey"`
	Amount uint
}
