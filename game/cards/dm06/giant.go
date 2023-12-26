package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func CantankerousGiant(c *match.Card) {

	c.Name = "Cantankerous Giant"
	c.Power = 8000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker)
}

func CliffcrushGiant(c *match.Card) {

	c.Name = "Cliffcrush Giant"
	c.Power = 7000
	c.Civ = civ.Nature
	c.Family = []string{family.Giant}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.AttackCreature); ok {

			if event.CardID != card.ID {
				return
			}

			if len(fx.FindFilter(card.Player, match.BATTLEZONE, func(card *match.Card) bool { return !card.Tapped && event.CardID != card.ID })) > 0 {
				ctx.InterruptFlow()
				ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack creatures", card.Name))
			}

		}

		if event, ok := ctx.Event.(*match.AttackPlayer); ok {

			if event.CardID != card.ID {
				return
			}

			if len(fx.FindFilter(card.Player, match.BATTLEZONE, func(card *match.Card) bool { return !card.Tapped && event.CardID != card.ID })) > 0 {
				ctx.InterruptFlow()
				ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack players", card.Name))
			}

		}
	})
}
