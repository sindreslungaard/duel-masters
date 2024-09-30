package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func BexTheOracle(c *match.Card) {

	c.Name = "Bex, the Oracle"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, fx.BlockerWhenNoShields))
}
