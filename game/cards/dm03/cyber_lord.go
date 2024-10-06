package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Emeral ...
func Emeral(c *match.Card) {

	c.Name = "Emeral"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.RotateShields(card, ctx, 1)
	}))

}

// Shtra ...
func Shtra(c *match.Card) {

	c.Name = "Shtra"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			"Shtra: Select 1 card from your manazone that will be sent to your hand",
			1,
			1,
			false,
		).Map(func(c *match.Card) {
			c.Player.MoveCard(c.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's manazone from their hand by Shtra", c.Name, c.Player.Username()))
		})

		ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")
		defer ctx.Match.EndWait(card.Player)

		fx.Select(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.MANAZONE,
			"Shtra: Select 1 card from your manazone that will be sent to your hand",
			1,
			1,
			false,
		).Map(func(c *match.Card) {
			c.Player.MoveCard(c.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved to %s's manazone from their hand by Shtra", c.Name, c.Player.Username()))
		})

	}))

}
