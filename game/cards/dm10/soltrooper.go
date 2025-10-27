package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// IkazTheSpydroid ...
func IkazTheSpydroid(c *match.Card) {

	c.Name = "Ikaz, The Spydroid"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.Soltrooper}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers,
		fx.When(fx.Blocks, func(card *match.Card, ctx *match.Context) {
			fx.Select(
				card.Player,
				ctx.Match,
				card.Player,
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: Choose one of your creatures in the battlezone. Untap it after the battle (%s is blocking).", card.Name, card.Name),
				1,
				1,
				false,
			).Map(func(x *match.Card) {
				ctx.ScheduleAfter(func() {
					if x.Zone != match.BATTLEZONE {
						return
					}

					x.Tapped = false
					ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was untapped by %s's effect.", x.Name, card.Name))
					ctx.Match.BroadcastState()
				})
			})
		}))

}
