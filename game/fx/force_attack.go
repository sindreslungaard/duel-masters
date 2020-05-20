package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// ForceAttack prevents the user from ending their turn if the card has not attacked this turn
func ForceAttack(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.EndTurnEvent); ok && card.Zone == match.BATTLEZONE {

		if ctx.Match.IsPlayerTurn(card.Player) && !card.HasCondition(cnd.SummoningSickness) && !card.Tapped {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s must attack before you can end your turn", card.Name))
			ctx.InterruptFlow()
		}

	}

}
