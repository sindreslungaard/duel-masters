package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Gigabolver ...
func Gigabolver(c *match.Card) {

	c.Name = "Gigabolver"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		fx.FilterShieldTriggers(ctx, func(x *match.Card) bool { return x.Civ != civ.Light })
	})

}
