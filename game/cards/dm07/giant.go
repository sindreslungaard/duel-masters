package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func HeadlongGiant(c *match.Card) {

	c.Name = "Headlong Giant"
	c.Power = 14000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 9
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Triplebreaker,
		fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {
			hand, err := card.Player.Container(match.HAND)
			if err != nil {
				return
			}

			if len(hand) == 0 {
				ctx.InterruptFlow()
				ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack if you don't have cards in hand", card.Name))
			}
		}),
		fx.When(fx.BlockerSelectionStep, fx.CantBeBlockedByPowerUpTo4000),
		fx.When(fx.AttackConfirmed, func(c *match.Card, ctx2 *match.Context) {
			fx.Select(c.Player, ctx2.Match, c.Player, match.HAND,
				fmt.Sprintf("%s: select a card to discard", c.Name), 1, 1, false,
			).Map(func(x *match.Card) {
				c.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD, c.ID)
			})
		}))

}
