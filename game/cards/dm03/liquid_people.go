package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AquaDeformer ...
func AquaDeformer(c *match.Card) {

	c.Name = "Aqua Deformer"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		cards := match.Search(card.Player, ctx.Match, card.Player, match.MANAZONE, "Aqua Deformer: Select 2 cards from your manazone that will be sent to your hand", 2, 2, false)

		for _, crd := range cards {
			card.Player.MoveCard(crd.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's hand from their mana zone", crd.Name, ctx.Match.PlayerRef(card.Player).Socket.User.Username))
		}

		ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")
		defer ctx.Match.EndWait(card.Player)

		opponentCards := match.Search(ctx.Match.Opponent(card.Player), ctx.Match, ctx.Match.Opponent(card.Player), match.MANAZONE, "Aqua Deformer: Select 2 cards from your manazone that will be sent to your hand", 2, 2, false)

		for _, crd := range opponentCards {
			ctx.Match.Opponent(card.Player).MoveCard(crd.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved to %s's hand from their mana zone", crd.Name, ctx.Match.PlayerRef(ctx.Match.Opponent(card.Player)).Socket.User.Username))
		}

	}))
}
