package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func LeQuistTheOracle(c *match.Card) {
	c.Name = "Le Quist, the Oracle"
	c.Power = 1500
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.WheneverThisAttacks(func(c *match.Card, ctx *match.Context) {
		filter := func(x *match.Card) bool { return x.Civ == civ.Fire || x.Civ == civ.Darkness }
		cards := make(map[string][]*match.Card)
		cards["Your creatures"] = fx.FindFilter(c.Player, match.BATTLEZONE, filter)
		cards["Opponent's creatures"] = fx.FindFilter(ctx.Match.Opponent(c.Player), match.BATTLEZONE, filter)

		fx.SelectMultipart(
			c.Player,
			ctx.Match,
			cards,
			"Le Quist, the Oracle: Select a card to tap or close to cancel",
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			x.Tapped = true
			fx.RemoveBlockerFromList(x, ctx)
		})

	}))

}
