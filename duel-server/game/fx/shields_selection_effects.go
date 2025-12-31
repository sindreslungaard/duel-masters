package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// ShieldsSelectionEffect adds the HasShieldsSelectionEffect condition in UntapStep
func ShieldsSelectionEffect(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.HasShieldsSelectionEffect, nil, nil)

	}

}
