package server

import "duel-masters/db"

// Message is the default message struct
type Message struct {
	Header string `json:"header"`
}

// DecksMessage lists the users decks
type DecksMessage struct {
	Header string          `json:"header"`
	Decks  []db.LegacyDeck `json:"decks"`
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
	CardID  string `json:"virtualId"`
	ImageID string `json:"uid"`
	Name    string `json:"name"`
	Civ     string `json:"civilization"`
	Flags   uint8  `json:"flags"`
}

// PlayerState stores information about the state of the current player
type PlayerState struct {
	Username   string         `json:"username"`
	Color      string         `json:"color"`
	Deck       int            `json:"deck"`
	HandCount  int            `json:"handCount"`
	Hand       []CardState    `json:"hand"`
	Shieldzone []string       `json:"shieldzone"`
	ShieldMap  map[string]int `json:"shieldMap"`
	Manazone   []CardState    `json:"manazone"`
	Graveyard  []CardState    `json:"graveyard"`
	Battlezone []CardState    `json:"playzone"`
}

// MatchState stores information about the current state of the match in the eyes of a given player
type MatchState struct {
	MyTurn       bool        `json:"myTurn"`
	HasAddedMana bool        `json:"hasAddedManaThisRound"`
	Me           PlayerState `json:"me"`
	Opponent     PlayerState `json:"opponent"`
	Spectator    bool        `json:"spectator"`
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
	Header            string      `json:"header"`
	ActionType        string      `json:"actionType"`
	Cards             []CardState `json:"cards"`
	Text              string      `json:"text"`
	MinSelections     int         `json:"minSelections"`
	MaxSelections     int         `json:"maxSelections"`
	Cancellable       bool        `json:"cancellable"`
	UnselectableCards []CardState `json:"unselectableCards"`
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

// LobbyChatMessage is used to store chat messages
type LobbyChatMessage struct {
	Username  string `json:"username"`
	Color     string `json:"color"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
	Removed   bool
}

// LobbyChatMessages is used to store chat messages
type LobbyChatMessages struct {
	Header   string             `json:"header"`
	Messages []LobbyChatMessage `json:"messages"`
}

// UserMessage holds information about users
type UserMessage struct {
	Username    string   `json:"username"`
	Color       string   `json:"color"`
	Hub         string   `json:"hub"`
	Permissions []string `json:"permissions"`
}

// UserListMessage is used to send a list of online users
type UserListMessage struct {
	Header string        `json:"header"`
	Users  []UserMessage `json:"users"`
}

// MatchMessage holds information about a match
type MatchMessage struct {
	ID          string `json:"id"`
	P1          string `json:"p1"`
	P1color     string `json:"p1color"`
	P2          string `json:"p2"`
	P2color     string `json:"p2color"`
	Name        string `json:"name"`
	Started     bool   `json:"spectate"`
	Matchmaking bool   `json:"matchmaking"`
}

// MatchesListMessage is used to list open matches
type MatchesListMessage struct {
	Header  string         `json:"header"`
	Matches []MatchMessage `json:"matches"`
}

type MatchRequestMessage struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	HostID     string `json:"host_id"`
	HostName   string `json:"host_name"`
	HostColor  string `json:"host_color"`
	GuestID    string `json:"guest_id"`
	GuestName  string `json:"guest_name"`
	GuestColor string `json:"guest_color"`
	Format     string `json:"format"`
	LinkCode   string `json:"link_code"`
}

type MatchReuestsListMessage struct {
	Header   string                `json:"header"`
	Requests []MatchRequestMessage `json:"requests"`
}

type MatchForwardMessage struct {
	Header string `json:"header"`
	ID     string `json:"id"`
}

// ShowCardsMessage is used to show the user n cards without an action to perform
type ShowCardsMessage struct {
	Header  string   `json:"header"`
	Message string   `json:"message"`
	Cards   []string `json:"cards"`
}

type PinnedMessages struct {
	Header   string   `json:"header"`
	Messages []string `json:"messages"`
}

type PlaySoundMessage struct {
	Header string `json:"header"`
	Sound  string `json:"sound"`
}
