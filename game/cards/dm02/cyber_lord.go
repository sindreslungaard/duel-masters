package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HypersquidWalter ...
func HypersquidWalter(c *match.Card) {

	c.Name = "Hypersquid Walter"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = family.CyberLord
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		ctx.ScheduleAfter(func() {
			card.Player.DrawCards(1)
		})

	}))

}
