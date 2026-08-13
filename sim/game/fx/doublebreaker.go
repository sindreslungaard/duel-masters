package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// Doublebreaker breaks two shields instead of 1 when attacking the player
func Doublebreaker(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.DoubleBreaker, true, card.ID)

	}

}
