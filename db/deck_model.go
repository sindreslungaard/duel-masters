package db

type LegacyDeck struct {
	UID      string   `json:"uid"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Public   bool     `json:"public"`
	Standard bool     `json:"standard"`
	Cards    []string `json:"cards"`
}

type Deck struct {
	UID      string `json:"uid"`
	Owner    string `json:"owner"`
	Name     string `json:"name"`
	Public   bool   `json:"public"`
	Standard bool   `json:"standard"`
	Cards    string `json:"cards"`
}

var Decks = conn().Collection("decks")
