package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// AmberPiercer ...
func AmberPiercer(c *match.Card) {

	c.Name = "Amber Piercer"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		ctx.ScheduleAfter(func() {

			fx.SelectFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.GRAVEYARD,
				"Amber Piercer: Select a card to return from your graveyard to your hand.",
				1,
				1,
				true,
				func(x *match.Card) bool { return x.HasCondition(cnd.Creature) },
			).Map(func(x *match.Card) {
				if x.ID == card.ID {
					return
				}
				card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s was returned to %s's hand from their graveyard by Amber Piercer", x.Name, x.Player.Username()))
			})

		})

	}))

}
