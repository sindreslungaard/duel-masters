package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// Adds the tap ability condition every turn
func TapAbility(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddUniqueSourceCondition(cnd.TapAbility, card.TapAbility, card.ID)
	}

}
