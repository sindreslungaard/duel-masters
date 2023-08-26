package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SpiralGrass ...
func SpiralGrass(c *match.Card) {

	c.Name = "Spiral Grass"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = []string{family.StarlightTree}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

			if event.Source == card && event.Blocked {
				card.Tapped = false
			}

		}

	})

}
