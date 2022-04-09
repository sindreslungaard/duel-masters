package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

func SpeedAttacker(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.CardMoved); ok {

		if event.CardID == card.ID && event.To == match.BATTLEZONE {

			ctx.ScheduleAfter(func() {
				card.RemoveCondition(cnd.SummoningSickness)
			})

		}

	}

}
