package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// Survivor adds the survivor condition every turn
func Survivor(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.Survivor, true, card.ID)

	}

}
