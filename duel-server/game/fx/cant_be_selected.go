package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// CantBeSelectedByOpp
func CantBeSelectedByOpp(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.CantBeSelectedByOpp, nil, card.ID)

	}

}
