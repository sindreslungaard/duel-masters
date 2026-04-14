package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AncientGiant ...
func AncientGiant(c *match.Card) {

	c.Name = "Ancient Giant"
	c.Power = 9000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.CantBeBlockedByDarkness)
}
