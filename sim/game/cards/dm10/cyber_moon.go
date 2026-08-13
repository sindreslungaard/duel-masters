package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ArdentLunatron ...
func ArdentLunatron(c *match.Card) {

	c.Name = "Ardent Lunatron"
	c.Power = 5000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberMoon}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers, fx.CantAttackCreatures, fx.BlockIfAbleWhenOppAttacks)

}

// HawkeyeLunatron ...
func HawkeyeLunatron(c *match.Card) {

	c.Name = "Hawkeye Lunatron"
	c.Power = 6000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.CyberMoon}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	// The searched card is not shown to the opponent.
	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, fx.SearchDeckTakeXCardsWithoutShowing(1)))

}
