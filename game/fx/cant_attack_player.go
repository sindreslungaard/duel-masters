package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// CantAttackPlayers prevents a card from attacking players
func CantAttackPlayers(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.AttackPlayer); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack players", card.Name))

		ctx.InterruptFlow()

	}

}
