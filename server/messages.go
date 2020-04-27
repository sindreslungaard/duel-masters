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
