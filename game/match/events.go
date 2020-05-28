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

// SpellCast is fired when a spell is cast, either from being played or from shield triggers
type SpellCast struct {
	CardID     string
	FromShield bool
}

// AttackPlayer is fired when the player attempts to use a creature to attack the player
type AttackPlayer struct {
	CardID   string
	Blockers []*Card
}

// AttackCreature is fired when the player attempts to use a creature to attack the player
type AttackCreature struct {
	CardID   string
	Blockers []*Card
}

// Battle is fired when two creatures are fighting, i.e. from attacking a creature or blocking an attack
type Battle struct {
	Attacker *Card
	Defender *Card
	Blocked  bool
}

// CreatureDestroyed is fired when a creature dies in battle or is destroyed from another source, such as a spell
type CreatureDestroyed struct {
	Card    *Card
	Source  *Card
	Blocked bool
}
