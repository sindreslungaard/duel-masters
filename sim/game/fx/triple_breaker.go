package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// Triplebreaker breaks three shields instead of 1 when attacking the player
func Triplebreaker(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.TripleBreaker, true, card.ID)

	}

}
