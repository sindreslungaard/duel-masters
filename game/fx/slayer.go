package fx

import (
	"duel-masters/game/match"
)

// Slayer destroys the source card when
func Slayer(card *match.Card, ctx *match.Context) {

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

		if event.Card == card {

			creature, err := ctx.Match.Opponent(card.Player).GetCard(event.Source.ID, match.BATTLEZONE)

			if err == nil {

				ctx.Match.Destroy(creature, card)

			}

		}

	}

}
