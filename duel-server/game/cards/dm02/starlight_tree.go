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

	c.Use(fx.Creature, fx.Blocker(), func(card *match.Card, ctx *match.Context) {
		if event, ok := ctx.Event.(*match.Battle); ok {
			if event.Defender == card && event.Blocked {
				ctx.ScheduleAfter(func() {
					if card.Zone == match.BATTLEZONE {
						card.Tapped = false
						ctx.Match.BroadcastState()
					}
				})
			}
		}
	})

}
