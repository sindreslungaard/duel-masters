package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

type SlayerCondition func(target *match.Card) bool

// Slayer destroys the source card when the card is destroyed
func Slayer(card *match.Card, ctx *match.Context) {
	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddCondition(cnd.Slayer, nil, card.ID)
	}
}

func ConditionalSlayer(condition SlayerCondition) func(card *match.Card, ctx *match.Context) {
	return func(card *match.Card, ctx *match.Context) {
		if _, ok := ctx.Event.(*match.UntapStep); ok {
			card.AddCondition(cnd.Slayer, condition, card.ID)
		}
	}
}

// Suicide destroys the card when it wins a battle
func Suicide(card *match.Card, ctx *match.Context) {
	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddCondition(cnd.DestroyAfterBattle, nil, card.ID)
	}
}
