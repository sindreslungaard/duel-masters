package fx

import (
	"duel-masters/game/events"
	"duel-masters/game/match"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, c *match.Context) {

	if _, ok := c.Event.(events.UntapStep); ok {

		if card.PlayerID == c.Match.PlayerTurnID {
			card.Tapped = false
		}

	}

}

func addToBattlezone(c *match.Context) {

}
