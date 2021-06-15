package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// PurplePiercer ...
func PurplePiercer(c *match.Card) {

	c.Name = "Purple Piercer"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = family.BrainJacker
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		//TODO cant be blocked by light, cant be attacked by light

	})

}
