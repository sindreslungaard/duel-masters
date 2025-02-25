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
		mana := fx.Find(card.Player, match.MANAZONE)

		if len(mana) <= 2 {
			mana.Map(func(x *match.Card) {
				x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to %s's hand from their mana zone", x.Name, ctx.Match.PlayerRef(x.Player).Socket.User.Username))
			})
		} else {
			fx.Select(card.Player,
				ctx.Match, card.Player,
				match.MANAZONE,
				"Aqua Deformer: Select 2 cards from your manazone that will be sent to your hand",
				2,
				2,
				false,
			).Map(func(x *match.Card) {
				x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to %s's hand from their mana zone", x.Name, ctx.Match.PlayerRef(x.Player).Socket.User.Username))
			})
		}

		ctx.Match.Wait(card.Player, "Waiting for your opponent to make an action")
		defer ctx.Match.EndWait(card.Player)

		opponentMana := fx.Find(ctx.Match.Opponent(card.Player), match.MANAZONE)

		if len(opponentMana) <= 2 {
			opponentMana.Map(func(x *match.Card) {
				x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to %s's hand from their mana zone", x.Name, ctx.Match.PlayerRef(x.Player).Socket.User.Username))
			})
		} else {
			fx.Select(
				ctx.Match.Opponent(card.Player),
				ctx.Match, ctx.Match.Opponent(card.Player),
				match.MANAZONE,
				"Aqua Deformer: Select 2 cards from your manazone that will be sent to your hand",
				2,
				2,
				false,
			).Map(func(x *match.Card) {
				x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
				ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to %s's hand from their mana zone", x.Name, ctx.Match.PlayerRef(x.Player).Socket.User.Username))
			})
		}
	}))
}
