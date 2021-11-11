package match

// EndTurnEvent is fired when a player attempts to end their turn
type EndTurnEvent struct {
}

// PlayCardEvent is fired when the player attempts to play a card
type PlayCardEvent struct {
	CardID string
}

// CardPlayedEvent is fired right before moving the card to the battle/spell zone, after choosing mana to tap
type CardPlayedEvent struct {
	CardID string
}

// ChargeManaEvent is fired when the player attempts to charge mana
type ChargeManaEvent struct {
	CardID string
}

// BrokenShieldEvent is fired right after a shield was broken
type BrokenShieldEvent struct {
	CardID string
}

// ShieldTriggerEvent is fired when a shield with shieldtrigger is broken
// can be cancelled to prevent the player from playing the card immediately
type ShieldTriggerEvent struct {
	Card *Card
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
	CardID              string
	Blockers            []*Card
	AttackableCreatures []*Card // list of cards that can be attacked
}

// AttackConfirmed is fired when either attacking a player or creature, but **after** the attack is validated
// to happen (does not fire when the attack is cancelled)
type AttackConfirmed struct {
	CardID   string
	Player   bool
	Creature bool
}

// Battle is fired when two creatures are fighting, i.e. from attacking a creature or blocking an attack
type Battle struct {
	Attacker *Card
	Defender *Card
	Blocked  bool
}

type CreatureDestroyedContext int

const (
	DestroyedInBattle = iota
	DestroyedBySpell
	DestroyedBySlayer
	DestroyedByMiscAbility
)

// CreatureDestroyed is fired when a creature dies in battle or is destroyed from another source, such as a spell
type CreatureDestroyed struct {
	Card    *Card
	Source  *Card
	Blocked bool
	Context CreatureDestroyedContext
}

// GetPowerEvent is fired whenever a card's power is to be used
type GetPowerEvent struct {
	Card      *Card
	Attacking bool
	Power     int
}
