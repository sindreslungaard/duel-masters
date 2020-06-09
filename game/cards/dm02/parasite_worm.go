package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ChaosWorm ...
func ChaosWorm(c *match.Card) {

	c.Name = "Chaos Worm"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = family.ParasiteWorm
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Chaos Worm: Select a creature from the opponent's battle zone and destroy it",
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card)
		})

	}))

}
