package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// ShieldTrigger returns the card to the players hand instead of the graveyard
func ShieldTrigger(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.ShieldTrigger, nil, nil)

	}

}
