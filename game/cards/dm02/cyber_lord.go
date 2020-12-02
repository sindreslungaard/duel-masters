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

// Corile ...
func Corile(c *match.Card) {

	c.Name = "Corile"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = family.CyberLord
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Corile: Move 1 of your opponent's creatures from their battlezone to the top of their deck",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			ctx.Match.MoveCard(x, match.DECK, card)
		})

	}))

}
