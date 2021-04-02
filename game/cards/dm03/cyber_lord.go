package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Emeral ...
func Emeral(c *match.Card) {

	c.Name = "Emeral"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = family.CyberLord
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID == card.ID && event.To == match.BATTLEZONE {

				nrShields, err := card.Player.Container(match.SHIELDZONE)

				if err != nil {
					return
				}

				if len(nrShields) < 1 {
					return
				}

				toShield := match.Search(card.Player, ctx.Match, card.Player, match.HAND, "Emeral: You may select 1 card from your hand and put it into the shield zone", 0, 1, true)

				if len(toShield) < 1 {
					return
				}

				toHand := fx.SelectBackside(
					card.Player,
					ctx.Match,
					card.Player,
					match.SHIELDZONE,
					"Emeral: Select 1 of your shields that will be moved to your hand",
					1,
					1,
					false,
				)

				for _, card := range toShield {
					card.Player.MoveCard(card.ID, match.HAND, match.SHIELDZONE)
				}

				for _, card := range toHand {
					card.Player.MoveCard(card.ID, match.SHIELDZONE, match.HAND)
				}

			}

		}

	})

}