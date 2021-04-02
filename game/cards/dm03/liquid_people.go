package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AquaDeformer ...
func AquaDeformer(c *match.Card) {

	c.Name = "Aqua Deformer"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = family.LiquidPeople
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				cards := match.Search(card.Player, ctx.Match, card.Player, match.MANAZONE, "Aqua Deformer: Select 2 cards from your manazone that will be sent to your hand", 2, 2, false)

				for _, card := range cards {
					card.Player.MoveCard(card.ID, match.MANAZONE, match.HAND)
				}

				ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")
				defer ctx.Match.EndWait(card.Player)

				opponentCards := match.Search(ctx.Match.Opponent(card.Player), ctx.Match, ctx.Match.Opponent(card.Player), match.MANAZONE, "Aqua Deformer: Select 2 cards from your manazone that will be sent to your hand", 2, 2, false)

				for _, card := range opponentCards {
					card.Player.MoveCard(card.ID, match.MANAZONE, match.HAND)
				}

			}

		}

	})
}