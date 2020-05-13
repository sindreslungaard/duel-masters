package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AquaHulcus ...
func AquaHulcus(c *match.Card) {

	c.Name = "Aqua Hulcus"
	c.Civ = civ.Water
	c.Family = family.LiquidPeople
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Draw1)

}
