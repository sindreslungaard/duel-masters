package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func Garatyano(c *match.Card) {

	c.Name = "Garatyano"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
		cardsOnTop := card.Player.PeekDeck(3)
		sortedCards := fx.OrderCards(
			card.Player,
			ctx.Match,
			cardsOnTop,
			"Order cards from top of your deck (First one will be at the bottom)",
		)

		for _, cID := range sortedCards {
			card.Player.MoveCardToFront(cID, match.DECK, match.DECK)
		}
	}

	c.Use(fx.Creature, fx.TapAbility)
}

func Biancus(c *match.Card) {

	c.Name = "Biancus"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = fx.GiveOwnCreatureCantBeBlocked

	c.Use(fx.Creature, fx.Blocker, fx.TapAbility)
}
