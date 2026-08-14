package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// AttackUntapped allows the card to attack untapped creatures
func AttackUntapped(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.AttackUntapped, true, card.ID)

	}

}
