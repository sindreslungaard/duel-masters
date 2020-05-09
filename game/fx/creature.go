package fx

import (
	"duel-masters/game/match"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, c *match.Context) {

	// Untap the card at the UntapStep
	if _, ok := c.Event.(*match.UntapStep); ok {

		if c.Match.IsPlayerTurn(card.Player) {
			card.Tapped = false
		}

	}

	// Move to the battlezone
	if event, ok := c.Event.(*match.PlayCardEvent); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		// OK, let's move me to the battlezone
		card.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE)

	}

}
