package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// NastashaChannelerOfSuns ...
func NastashaChannelerOfSuns(c *match.Card) {
	c.Name = "Nastasha, Channeler of Suns"
	c.Power = 6000
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.BreakShield, func(card *match.Card, ctx *match.Context) {
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
				"Those shields are about to be broken. You may choose a shield to protect and destroy Nastasha, Channeler of Suns.",
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

				ctx.Match.ReportActionInChat(card.Player, "Nastasha, Channeler of Suns was destoryed to protect a shield.")
				ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
			})
		})
	}))
}
