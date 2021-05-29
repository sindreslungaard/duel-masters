package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// Blocker adds the card to a list of blockers when a creature/player is attacked
func Blocker(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.Blocker, true, card.ID)

	}

}
