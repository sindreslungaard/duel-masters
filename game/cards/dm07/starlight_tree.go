package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func PulsarTree(c *match.Card) {

	c.Name = "Pulsar Tree"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.StarlightTree}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.BreakShield, func(card *match.Card, ctx *match.Context) {

		event, ok := ctx.Event.(*match.BreakShieldEvent)
		if !ok {
			return
		}

		if len(event.Cards) < 1 || event.Cards[0].Player != card.Player {
			return
		}

		// Adding ScheduleAfter so it works well with cards that interrupt the event
		ctx.ScheduleAfter(func() {
			fx.SelectBacksideFilter(card.Player, ctx.Match, card.Player, match.SHIELDZONE,
				"Those shields are about to be broken. You may choose a shield to protect and destroy Pulsar Tree.",
				1, 1, true, func(c *match.Card) bool {
					for _, shield := range event.Cards {
						if shield == c {
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

				ctx.Match.ReportActionInChat(card.Player, "Pulsar tree destoryed to protect a shield.")
				ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
			})
		})
	}))
}
