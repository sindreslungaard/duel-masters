package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Pokolul ...
func Pokolul(c *match.Card) {

	c.Name = "Pokolul"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	isAttacking := false

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {

			if event.CardID == card.ID {
				isAttacking = true
			} else {
				isAttacking = false
			}

		}

		if _, ok := ctx.Event.(*match.ShieldTriggerEvent); ok && isAttacking {
			card.Tapped = false
			ctx.Match.Chat("Server", fmt.Sprintf("%s untapped", card.Name))
		}

	})

}
