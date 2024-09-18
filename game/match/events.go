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

// BreakShieldEvent is fired before the shield is broken
type BreakShieldEvent struct {
	CardID string
	Source string // the card id that caused the shield to break
}

// BrokenShieldEvent is fired right after a shield was broken
type BrokenShieldEvent struct {
	CardID string
	Source string // the card id that caused the shield to break
}

// ShieldTriggerEvent is fired when a shield with shieldtrigger is broken
// can be cancelled to prevent the player from playing the card immediately
type ShieldTriggerEvent struct {
	Cards           []*Card
	UnplayableCards []*Card
	Source          string // the card id that caused the shield to break
}

// ShieldTriggerPlayedEvent is fired when a shield trigger is played
type ShieldTriggerPlayedEvent struct {
	Card   *Card
	Source string // the card id that caused the shield to break
}

// MoveCard is fired from the *Player.MoveCard method before moving a card between containers
type MoveCard struct {
	CardID string
	From   string
	To     string
	Source string // What caused the card to move, usually the ID of a card
}

// CardMoved is fired from the *Player.MoveCard method after moving a card between containers
type CardMoved struct {
	CardID        string
	From          string
	To            string
	Source        string // What caused the card to move, usually the ID of a card
	MatchPlayerID byte
}

// SpellCast is fired when a spell is cast, either from being played or from shield triggers
type SpellCast struct {
	CardID        string
	FromShield    bool
	MatchPlayerID byte
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
	Attacker      *Card
	AttackerPower int
	Defender      *Card
	DefenderPower int
	Blocked       bool
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

// TapAbility is fired when the user activates a cards tap ability instead of attacking
type TapAbility struct {
	CardID string
}
