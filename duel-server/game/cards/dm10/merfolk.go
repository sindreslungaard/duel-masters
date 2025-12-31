package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TidePatroller ...
func TidePatroller(c *match.Card) {

	c.Name = "Tide Patroller"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.Merfolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker())

}
