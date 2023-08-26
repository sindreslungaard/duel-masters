package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Gigastand ...
func Gigastand(c *match.Card) {

	c.Name = "Gigastand"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.CardMoved); ok {

			if event.CardID != card.ID || event.To != match.GRAVEYARD || event.From != match.BATTLEZONE {
				return
			}

			hand := match.Search(
				card.Player,
				ctx.Match,
				card.Player,
				match.HAND,
				"Gigastand was destroyed, you may return it to your hand by discarding another card from your hand. Close to cancel.",
				1,
				1,
				true,
			)

			if len(hand) < 1 {
				return
			}

			card.Player.MoveCard(card.ID, match.GRAVEYARD, match.HAND)
			hand[0].Player.MoveCard(hand[0].ID, match.HAND, match.GRAVEYARD)

			ctx.Match.Chat("Server", fmt.Sprintf("Gigastand was moved to %s's hand instead of the graveyard", card.Player.Username()))
			ctx.Match.Chat("Server", fmt.Sprintf("%s was discarded from %s's hand", hand[0].Name, hand[0].Player.Username()))

		}

	})

}
