package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// CantBeBlocked allows the card to attack without being blocked
func CantBeBlocked(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.CantBeBlocked, nil, card.ID)

	}

}
