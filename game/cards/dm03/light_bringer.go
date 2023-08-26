package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AlessTheOracle ...
func AlessTheOracle(c *match.Card) {

	c.Name = "Aless, the Oracle"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.ReturnToShield)
}
