package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ProwlingElephish ...
func ProwlingElephish(c *match.Card) {

	c.Name = "Prowling Elephish"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker())

}
