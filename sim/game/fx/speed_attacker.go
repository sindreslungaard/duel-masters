package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

func SpeedAttacker(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.SpeedAttacker, true, card.ID)

	}

}
