package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KolonTheOracle ...
func KolonTheOracle(c *match.Card) {

	c.Name = "Kolon, the Oracle"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.When(fx.Summoned, fx.MayTapOpCreature))
}
