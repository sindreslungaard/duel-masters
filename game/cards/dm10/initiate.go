package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// GlaisMejiculaTheExtreme ...
func GlaisMejiculaTheExtreme(c *match.Card) {

	c.Name = "Glais Mejicula, the Extreme"
	c.Power = 5500
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.BreakShield, func(card *match.Card, ctx *match.Context) {

		event, ok := ctx.Event.(*match.BreakShieldEvent)
		if !ok {
			return
		}

		if len(event.Cards) < 1 || event.Cards[0].Player != card.Player {
			return
		}

		cardsInHand, err := card.Player.Container(match.HAND)

		if err != nil || len(cardsInHand) < 2 {
			return
		}

		maxShields := min(len(event.Cards), len(cardsInHand)/2)

		// Adding ScheduleAfter so it works well with cards that interrupt the event
		ctx.ScheduleAfter(func() {
			fx.SelectBacksideFilter(card.Player, ctx.Match, card.Player, match.SHIELDZONE,
				fmt.Sprintf("Those shields are about to be broken. You may choose up to %v shields to protect and discard 2 cards from your hand instead, for each shield chosen.", maxShields),
				1, maxShields, true, func(x *match.Card) bool {
					for _, shield := range event.Cards {
						if shield == x {
							return true
						}
					}
					return false
				},
			).Map(func(x *match.Card) {
				if len(event.Cards) < 2 {
					ctx.InterruptFlow()
				} else {
					var validShields []*match.Card
					for _, shield := range event.Cards {
						if shield != x {
							validShields = append(validShields, shield)
						}
					}
					event.Cards = validShields
				}

				fx.Select(
					card.Player,
					ctx.Match,
					card.Player,
					match.HAND,
					"Select 2 cards from your hand to discard.",
					2,
					2,
					false,
				).Map(func(x *match.Card) {
					_, err := x.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)

					if err == nil {
						ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s discarded %s from his hand to protect a shield.", card.Player.Username(), x.Name))
					}
				})
			})
		})
	}))

}
