package fx

import (
	"duel-masters/game/events"
	"duel-masters/game/match"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, c *match.Context) {

	// Untap the card at the UntapStep
	if _, ok := c.Event.(*events.UntapStep); ok {

		if c.Match.IsPlayerTurn(card.Player) {
			card.Tapped = false
		}

	}

	// Check for and tap required mana when played, move to the battlezone
	if event, ok := c.Event.(*events.PlayCard); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		// OK, but do my player have the necessary untapped mana to play me?
		/* untappedMana := make([]*match.Card, 0)
		for _, card := range c.Match.CurrentPlayer().Player {

		} */

		// OK, let's move me to the battlezone
		card.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE)

	}

}
