package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Marinomancer ...
func Marinomancer(c *match.Card) {

	c.Name = "Marinomancer"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned,  func(card *match.Card, ctx *match.Context) {

		cards := card.Player.PeekDeck(3)

		for _, toMove := range cards {

			if toMove.Civ == civ.Light || toMove.Civ == civ.Darkness {
				card.Player.MoveCard(toMove.ID, match.DECK, match.HAND)
				ctx.Match.Chat("Server", fmt.Sprintf("%s put %s into the hand from the top of their deck", card.Player.Username(), toMove.Name))
			} else {
				card.Player.MoveCard(toMove.ID, match.DECK, match.GRAVEYARD)
				ctx.Match.Chat("Server", fmt.Sprintf("%s put %s into the graveyard from the top of their deck", card.Player.Username(), toMove.Name))
			}
		}

	}))

}