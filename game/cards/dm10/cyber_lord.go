package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Zaltan ...
func Zaltan(c *match.Card) {

	c.Name = "Zaltan"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.AnotherOwnCyberVirusSummoned, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			fmt.Sprintf("%s's effect: You may discard up to 2 cards from your hand. For each card you discard, choose a creature in the battle zone and return it to its owner's hand.", card.Name),
			1,
			2,
			true,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved from %s's hand to his graveyard by %s", x.Name, x.Player.Username(), card.Name))

			fx.ReturnCreatureToOwnersHand(card, ctx)
		})
	}))

}
