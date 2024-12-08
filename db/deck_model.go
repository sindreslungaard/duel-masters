package db

type LegacyDeck struct {
	UID      string   `json:"uid"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Public   bool     `json:"public"`
	Standard bool     `json:"standard"`
	Cards    []string `json:"cards"`
	Event    string   `json:"event"`
	Cardpool []string `json:"cardpool"`
}

type Deck struct {
	UID      string `json:"uid"`
	Owner    string `json:"owner"`
	Name     string `json:"name"`
	Public   bool   `json:"public"`
	Standard bool   `json:"standard"`
	Cards    string `json:"cards"`
	Event    string `json:"event"`
	Cardpool string `json:"cardpool"`
}

var Decks = conn().Collection("decks")
