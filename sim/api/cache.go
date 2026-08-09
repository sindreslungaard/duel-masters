package api

import (
	"duel-masters/game/cards"
	"duel-masters/game/match"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// CardInfo describes the card catalog exposed by the simulator API.
type CardInfo struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
	// Civilizations lists every civilization printed on the card. Multicolored
	// cards have more than one and count as a card of each of them.
	Civilizations []string `json:"civilizations"`
	// ManaRequirement lists the civilizations that must all be present among the
	// mana paid for the card, one card of each.
	ManaRequirement []string `json:"manaRequirement"`
	Family          []string `json:"family"`
	ManaCost        int      `json:"manaCost"`
	Set             string   `json:"set"`
	Type            string   `json:"type"`
	Text            string   `json:"text"`
}

// Register holds all the card info
var register = make([]CardInfo, 0)
var mutex = &sync.Mutex{}

// CreateCardCache loads all cards and creates a cache of the static data
func CreateCardCache() {

	cardsFromJsonMap := readFromJson()

	for setID, set := range cards.Sets {

		for uid, c := range *set {

			if c == nil {
				continue
			}

			card := &match.Card{}

			c(card)

			entry := CardInfo{
				UID:             uid,
				Name:            card.Name,
				Civilizations:   card.Civs,
				ManaRequirement: card.ManaRequirement,
				Set:             setID,
				ManaCost:        card.ManaCost,
				Type:            "Creature",
			}

			if len(card.Family) > 0 {
				entry.Family = card.Family
			} else {
				entry.Type = "Spell"
			}

			if _, ok := cardsFromJsonMap[card.Name]; ok {
				entry.Text = cardsFromJsonMap[card.Name].Text
			} else {
				logrus.Warnf("Card '%s' not found in json file", card.Name)
			}

			register = append(register, entry)

		}

	}

	logrus.Infof("Loaded %v cards into the cache from %v sets", len(register), len(cards.Sets))

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

type CardsFromJson struct {
	Cards []CardFromJson `json:"cards"`
}

type CardFromJson struct {
	Civilizations []string            `json:"civilizations"`
	ManaCost      int                 `json:"cost"`
	Name          string              `json:"name"`
	Power         string              `json:"power"`
	Printings     []PrintingsFromJson `json:"printings"`
	Subtypes      []string            `json:"subtypes"`
	Supertypes    []string            `json:"supertypes"`
	Text          string              `json:"text"`
	Type          string              `json:"type"`
}

type PrintingsFromJson struct {
	Set         string `json:"set"`
	Id          string `json:"id"`
	Rarity      string `json:"rarity"`
	Flavor      string `json:"flavor"`
	Illustrator string `json:"illustrator"`
}

func readFromJson() map[string]CardFromJson {
	jsonFileName := "./sim/DuelMastersCards.json"
	jsonFile, err := os.Open(jsonFileName)
	if err != nil {
		logrus.Error(fmt.Sprintf("Error loading %s", jsonFileName), err)
		return nil
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var cards CardsFromJson
	json.Unmarshal(byteValue, &cards)

	logrus.Infof("Loaded %v card details from %s", len(cards.Cards), jsonFileName)

	cardsMap := make(map[string]CardFromJson)

	for _, card := range cards.Cards {
		cardsMap[card.Name] = card
	}

	return cardsMap
}
