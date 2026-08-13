package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// Evolve into any family allows the card to evolve into any other family
func EvolveIntoAnyFamily(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.EvolveIntoAnyFamily, true, card.ID)

	}

}
