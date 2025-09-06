package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// StratosphereGiant ...
func StratosphereGiant(c *match.Card) {

	c.Name = "Stratosphere Giant"
	c.Power = 13000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Triplebreaker, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.HAND,
			"Choose up to 2 creatures in your hand and put them into the battlezone.",
			1,
			2,
			true,
			func(x *match.Card) bool {
				return fx.CanBeSummoned(ctx.Match.Opponent(card.Player), x)
			},
			false,
		).Map(func(x *match.Card) {
			fx.ForcePutCreatureIntoBZ(ctx, x, match.HAND, card)
		})
	}))
}
