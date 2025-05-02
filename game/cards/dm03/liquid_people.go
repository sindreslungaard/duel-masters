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
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.MANAZONE,
			fmt.Sprintf("%s: Select 2 cards from your manazone that will be sent to your hand", card.Name),
			2,
			2,
			false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's hand from their mana zone", x.Name, ctx.Match.PlayerRef(card.Player).Socket.User.Username))
		})

		fx.Select(
			ctx.Match.Opponent(card.Player),
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.MANAZONE, fmt.Sprintf("%s: Select 2 cards from your manazone that will be sent to your hand", card.Name),
			2,
			2,
			false,
		).Map(func(x *match.Card) {
			x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(x.Player, fmt.Sprintf("%s was moved to %s's hand from their mana zone", x.Name, ctx.Match.PlayerRef(ctx.Match.Opponent(card.Player)).Socket.User.Username))
		})
	}))
}
