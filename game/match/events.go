package match

// EndTurnEvent is fired when a player attempts to end their turn
type EndTurnEvent struct {
}

// PlayCardEvent is fired when the player attempts to play a card
type PlayCardEvent struct {
	CardID string
}

// ChargeManaEvent is fired when the player attempts to charge mana
type ChargeManaEvent struct {
	CardID string
}

// CardMoved is fired from the *Player.MoveCard method after moving a card between containers
type CardMoved struct {
	CardID string
	From   string
	To     string
}

// AttackPlayer is fired when the player attempts to use a creature to attack the player
type AttackPlayer struct {
	CardID   string
	Blockers []*Card
}
