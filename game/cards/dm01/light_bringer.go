package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ChiliasTheOracle ...
func ChiliasTheOracle(c *match.Card) {

	c.Name = "Chilias, the Oracle"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = family.BeastFolk
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.ReturnToHand)

}
