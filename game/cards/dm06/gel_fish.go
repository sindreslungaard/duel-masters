package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func RaptorFish(c *match.Card) {

	c.Name = "Raptor Fish"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		hand := fx.Find(card.Player, match.HAND)
		toDraw := len(hand)

		if toDraw < 1 {
			return
		}

		for _, x := range hand {
			x.Player.MoveCard(x.ID, match.HAND, match.DECK)
		}

		card.Player.ShuffleDeck()
		card.Player.DrawCards(toDraw)
	}))

}
