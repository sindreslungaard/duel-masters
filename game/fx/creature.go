package fx

import (
	"duel-masters/game/events"
	"duel-masters/game/match"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, c *match.Context) {

	if _, ok := c.Event.(events.UntapStep); ok {

		if c.Match.PlayerTurn() == card.Player {
			card.Tapped = false
		}

	}

	if _, ok := c.Event.(events.PlayCardEvent); ok {

		if c.Match.PlayerTurn() == card.Player {
			card.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE)
		}

	}

}
