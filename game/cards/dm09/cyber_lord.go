package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// EmperorMaroll ...
func EmperorMaroll(c *match.Card) {

	c.Name = "Emperor Maroll"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Evolution,
		fx.When(fx.AnotherOwnCreatureSummoned, func(card *match.Card, ctx *match.Context) {

			if card.Zone != match.BATTLEZONE {
				return
			}

			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to the hand", card.Name))

		}),
		func(card *match.Card, ctx *match.Context) {

			if event, ok := ctx.Event.(*match.Battle); ok {
				if event.Blocked && event.Attacker == card {
					ctx.InterruptFlow()
					card.Tapped = true

					_, err := event.Defender.Player.MoveCard(event.Defender.ID, match.BATTLEZONE, match.HAND, card.ID)
					if err != nil {
						return
					}

					ctx.Match.ReportActionInChat(event.Defender.Player, fmt.Sprintf("%s was returned to its owner's hand instead of blocking due to %s's effect.", event.Defender.Name, card.Name))
				}
			}

		})
}
