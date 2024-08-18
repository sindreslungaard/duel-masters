package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// LagunaLightningEnforcer ...
func LagunaLightningEnforcer(c *match.Card) {

	c.Name = "Laguna, Lightning Enforcer"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = []string{family.Berserker}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		ctx.ScheduleAfter(func() {

			fx.SelectFilterFullList(card.Player,
				ctx.Match,
				card.Player,
				match.DECK,
				"Choose a spell from your deck that will be shown to your opponent",
				1,
				1,
				true,
				func(x *match.Card) bool { return x.HasCondition(cnd.Spell) },
				true,
			).Map(func(x *match.Card) {
				x.Player.MoveCard(x.ID, match.DECK, match.HAND, card.ID)
				ctx.Match.Chat("Server", fmt.Sprintf("%s retrieved %s from their deck to their hand", x.Player.Username(), x.Name))
			})

			card.Player.ShuffleDeck()
			ctx.Match.Chat("Server", fmt.Sprintf("%s shuffled their deck", card.Player.Username()))

		})

	}))

}
