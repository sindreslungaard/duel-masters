package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DomeShell ...
func DomeShell(c *match.Card) {

	c.Name = "Dome Shell"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = family.ColonyBeetle
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}
