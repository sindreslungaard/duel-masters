package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// MishaChannelerOfSuns ...
func MishaChannelerOfSuns(c *match.Card) {

	c.Name = "Misha, Channeler of Suns"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.CantBeAttackedIf(func(attacker *match.Card) bool {
		return attacker.SharesAFamily(family.Dragons)
	}))
}

// SashaChannelerOfSuns ...
func SashaChannelerOfSuns(c *match.Card) {

	c.Name = "Sasha, Channeler of Suns"
	c.Power = 9500
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker, fx.DragonBlocker(), func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.Battle); ok {
			if event.Attacker == card && event.Defender.SharesAFamily(family.Dragons) {
				event.AttackerPower += 6000
			} else if event.Defender == card && event.Attacker.SharesAFamily(family.Dragons) {
				event.DefenderPower += 6000
			}
		}

	})

}

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
			fx.SelectBacksideFilter(
				card.Player,
				ctx.Match,
				card.Player,
				match.SHIELDZONE,
				fmt.Sprintf("Those shields are about to be broken. You may choose a shield to protect and destroy %s instead.", card.Name),
				1,
				1,
				true,
				func(c *match.Card) bool {
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
					validShields := make([]*match.Card, 0)

					for _, shield := range event.Cards {
						if shield != x {
							validShields = append(validShields, shield)
						}
					}

					event.Cards = validShields
				}

				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was destoryed to protect a shield.", card.Name))
				ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
			})
		})
	}))
}
