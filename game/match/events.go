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
