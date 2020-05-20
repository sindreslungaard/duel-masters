package fx

import (
	"duel-masters/game/match"
)

// CantBeBlocked allows the card to attack without being blocked
func CantBeBlocked(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.AttackPlayer); ok {

		if event.CardID != card.ID {
			return
		}

		ctx.ScheduleAfter(func() {
			event.Blockers = make([]*match.Card, 0)
		})

	}

	if event, ok := ctx.Event.(*match.AttackCreature); ok {

		if event.CardID != card.ID {
			return
		}

		ctx.ScheduleAfter(func() {
			event.Blockers = make([]*match.Card, 0)
		})

	}

}
