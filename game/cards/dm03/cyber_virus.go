package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// StingerBall ...
func StingerBall(c *match.Card) {

	c.Name = "Stinger Ball"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		ctx.ScheduleAfter(func() {

			fx.SelectBackside(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.SHIELDZONE,
				"Stinger Ball: Select 1 of your opponent's shields that will be shown to you",
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

		})

	}))
}
