package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// FonchTheOracle ...
func FonchTheOracle(c *match.Card) {

	c.Name = "Fonch, the Oracle"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Fonch, the Oracle: Select a creature from the opponent's battle zone that will be tapped",
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			x.Tapped = true
		})

	}))

}

// WynTheOracle ...
func WynTheOracle(c *match.Card) {

	c.Name = "Wyn, the Oracle"
	c.Power = 1500
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		fx.SelectBackside(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.SHIELDZONE,
			"Wyn, the Oracle: Select 1 of your opponent's shields that will be shown to you",
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.ShowCards(
				card.Player,
				"Your opponent's shield:",
				[]string{x.ImageID},
			)
		})

	}))

}
