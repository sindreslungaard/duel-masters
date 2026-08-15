package server

// Message is the default message struct
type Message struct {
	Header string `json:"header"`
}

// DecksMessage lists the users decks
type DeckMessage struct {
	UID      string   `json:"uid"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Public   bool     `json:"public"`
	Standard bool     `json:"standard"`
	Cards    []string `json:"cards"`
}

type DecksMessage struct {
	Header string        `json:"header"`
	Decks  []DeckMessage `json:"decks"`
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

// ShieldState stores information about the state of a shield
type ShieldState struct {
	CardID  string `json:"virtualId"`
	ImageID string `json:"uid"`
	Flags   uint8  `json:"flags"`
}

// PlayerState stores information about the state of the current player
type PlayerState struct {
	Username   string         `json:"username"`
	Color      string         `json:"color"`
	Deck       int            `json:"deck"`
	HandCount  int            `json:"handCount"`
	Hand       []CardState    `json:"hand"`
	Shieldzone []ShieldState  `json:"shieldzone"`
	ShieldMap  map[string]int `json:"shieldMap"`
	Manazone   []CardState    `json:"manazone"`
	Graveyard  []CardState    `json:"graveyard"`
	Battlezone []CardState    `json:"playzone"`
}

// MatchState stores information about the current state of the match in the eyes of a given player
type MatchState struct {
	MyTurn        bool        `json:"myTurn"`
	HasAddedMana  bool        `json:"hasAddedManaThisRound"`
	HasAttacked   bool        `json:"hasAttackedThisRound"`
	CanChargeMana bool        `json:"canChargeManaThisRound"`
	Me            PlayerState `json:"me"`
	Opponent      PlayerState `json:"opponent"`
	Spectator     bool        `json:"spectator"`
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

type DuelFinishedPlayer struct {
	UID      string `json:"uid,omitempty"`
	Username string `json:"username,omitempty"`
}

type DuelFinishedMessage struct {
	Header               string              `json:"header"`
	DuelID               string              `json:"duelId"`
	Winner               *DuelFinishedPlayer `json:"winner,omitempty"`
	MatchResultGenerated bool                `json:"matchResultGenerated"`
	WonByDisconnect      bool                `json:"wonByDisconnect"`
	Turns                int                 `json:"turns"`
	DurationSeconds      int64               `json:"durationSeconds"`
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
	Choices           []string    `json:"choices"`
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

// ShowCardsMessage is used to show the user n cards without an action to perform
type ShowCardsMessage struct {
	Header  string   `json:"header"`
	Message string   `json:"message"`
	Cards   []string `json:"cards"`
}

type PlaySoundMessage struct {
	Header string `json:"header"`
	Sound  string `json:"sound"`
}
