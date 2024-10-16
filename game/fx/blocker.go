package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

type BlockerCondition func(target *match.Card) bool

// Blocker adds the card to a list of blockers when a creature/player is attacked
func Blocker(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.Blocker, true, card.ID)

	}

}

func ConditionalBlocker(condition BlockerCondition) func(card *match.Card, ctx *match.Context) {
	return func(card *match.Card, ctx *match.Context) {
		if _, ok := ctx.Event.(*match.UntapStep); ok {
			card.AddCondition(cnd.Blocker, condition, card.ID)
		}
	}
}
