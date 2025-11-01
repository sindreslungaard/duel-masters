package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TulkTheOracle ...
func TulkTheOracle(c *match.Card) {

	c.Name = "Tulk, the Oracle"
	c.Civ = civ.Light
	c.Power = 500
	c.Family = []string{family.LightBringer}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}
