package server

import "duel-masters/db"

// Message is the default message struct
type Message struct {
	Header string `json:"header"`
}

// DecksMessage lists the users decks
type DecksMessage struct {
	Header string    `json:"header"`
	Decks  []db.Deck `json:"decks"`
}

// ChatMessage stores information about a chat message
type ChatMessage struct {
	Header  string `json:"header"`
	Message string `json:"message"`
	Sender  string `json:"sender"`
	Color   string `json:"color"`
}

// CardState stores information about the state of a card
type CardState struct {
	CardID      string `json:"virtualId"`
	ImageID     string `json:"uid"`
	Name        string `json:"name"`
	Civ         string `json:"civilization"`
	Tapped      bool   `json:"tapped"`
	CanBePlayed bool   `json:"canBePlayed"`
}

// PlayerState stores information about the state of the current player
type PlayerState struct {
	Deck       int         `json:"deck"`
	Hand       []CardState `json:"hand"`
	Shieldzone []string    `json:"shieldzone"`
	Manazone   []CardState `json:"manazone"`
	Graveyard  []CardState `json:"graveyard"`
	Battlezone []CardState `json:"playzone"`
}

// MatchState stores information about the current state of the match in the eyes of a given player
type MatchState struct {
	MyTurn       bool        `json:"myTurn"`
	HasAddedMana bool        `json:"hasAddedManaThisRound"`
	Me           PlayerState `json:"me"`
	Opponent     PlayerState `json:"opponent"`
}

// MatchStateMessage is the message that should be sent to the client for state updates
type MatchStateMessage struct {
	Header string     `json:"header"`
	State  MatchState `json:"state"`
}

// WarningMessage is used to send a warning to a player
type WarningMessage struct {
	Header  string `json:"header"`
	Message string `json:"message"`
}

// ActionMessage is used to prompt the user to make a selection of the specified cards
type ActionMessage struct {
	Header        string      `json:"header"`
	Cards         []CardState `json:"cards"`
	Text          string      `json:"text"`
	MinSelections int         `json:"minSelections"`
	MaxSelections int         `json:"maxSelections"`
	Cancellable   bool        `json:"cancellable"`
}

// MultipartActionMessage is used to prompt the user to make a selection of the specified cards
type MultipartActionMessage struct {
	Header        string                 `json:"header"`
	Cards         map[string][]CardState `json:"cards"`
	Text          string                 `json:"text"`
	MinSelections int                    `json:"minSelections"`
	MaxSelections int                    `json:"maxSelections"`
	Cancellable   bool                   `json:"cancellable"`
}

// ActionWarningMessage is used to apply an error
type ActionWarningMessage struct {
	Header  string `json:"header"`
	Message string `json:"message"`
}

// WaitMessage is used to send a waiting popup with a message
type WaitMessage struct {
	Header  string `json:"header"`
	Message string `json:"message"`
}
