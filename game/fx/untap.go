package fx

import (
	"duel-masters/game/match"
)

// Untap untaps the card at each untap step, even the opponents
func Untap(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.EndOfTurnStep); ok {
		if card.Zone == match.BATTLEZONE {
			card.Tapped = false
		}
	}

}
