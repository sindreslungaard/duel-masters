package db

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Deck struct {
	gorm.Model
	ID       uint   `json:"uid"`
	UserID   uint   `json:"owner"`
	Name     string `json:"name"`
	Public   bool   `json:"public"`
	Standard bool   `json:"standard"`
}

type DeckCards struct {
	gorm.Model
	DeckID uint `gorm:"primaryKey"`
	CardID uint `gorm:"primaryKey"`
	Amount uint
}

type DeckAggregate struct {
	ID       uint     `json:"uid"`
	UserID   uint     `json:"owner"`
	Name     string   `json:"name"`
	Public   bool     `json:"public"`
	Standard bool     `json:"standard"`
	Cards    []string `json:"cards"`
}

func DecksForUserAggregated(conn *gorm.DB, userID uint) []DeckAggregate {
	var decks []Deck
	conn.Find(&decks, "user_id = ?", userID)

	var cards []DeckCards
	conn.Select("deck_cards.*").Joins("inner join deck_cards on deck_cards.deck_id = decks.id").Joins("inner join users on users.id = decks.user_id").Where("users.id = ?", userID).Find(&cards)

	decksMap := map[uint]*DeckAggregate{}

	for _, deck := range decks {
		decksMap[deck.ID] = &DeckAggregate{
			ID:       deck.ID,
			UserID:   deck.UserID,
			Name:     deck.Name,
			Public:   deck.Public,
			Standard: deck.Standard,
			Cards:    []string{},
		}
	}

	for _, card := range cards {
		d, ok := decksMap[card.DeckID]
		if !ok {
			logrus.Warn("Skipped card for deck with id ", card.DeckID)
			continue
		}

		d.Cards = append(d.Cards, fmt.Sprint(card.ID))
	}

	aggregare := []DeckAggregate{}

	for _, d := range decksMap {
		aggregare = append(aggregare, *d)
	}

	return aggregare
}
