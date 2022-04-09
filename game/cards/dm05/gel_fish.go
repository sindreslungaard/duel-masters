package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SeaSlug ...
func SeaSlug(c *match.Card) {

	c.Name = "Sea Slug"
	c.Power = 6000
	c.Civ = civ.Water
	c.Family = family.GelFish
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantBeBlocked)

}
