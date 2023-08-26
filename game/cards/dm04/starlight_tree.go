package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AmberGrass ...
func AmberGrass(c *match.Card) {

	c.Name = "Amber Grass"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.StarlightTree}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.ShieldTrigger)
}
