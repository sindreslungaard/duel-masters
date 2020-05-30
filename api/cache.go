package api

import (
	"duel-masters/game/cards"
	"duel-masters/game/match"
	"sync"

	"github.com/sirupsen/logrus"
)

// CardInfo struct is used for the card database api
type CardInfo struct {
	UID          string `json:"uid"`
	Name         string `json:"name"`
	Civilization string `json:"civilization"`
	Set          string `json:"set"`
}

// Register holds all the card info
var register []CardInfo = make([]CardInfo, 0)
var mutex *sync.Mutex = &sync.Mutex{}

// CreateCardCache loads all cards and creates a cache of the static data
func CreateCardCache() {

	for uid, c := range cards.DM01 {

		card := &match.Card{}

		c(card)

		register = append(register, CardInfo{
			UID:          uid,
			Name:         card.Name,
			Civilization: card.Civ,
			Set:          "dm-01",
		})

	}

	logrus.Infof("Loaded %v cards into the cache", len(register))

}

// GetCache returns a copy of the cache
func GetCache() []CardInfo {
	return register
}

// CacheHas returns true if the specified uid exist in the cache
func CacheHas(uid string) bool {

	mutex.Lock()

	defer mutex.Unlock()

	for _, c := range register {
		if c.UID == uid {
			return true
		}
	}

	return false

}
