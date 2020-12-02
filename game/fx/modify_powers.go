package fx

import "duel-masters/game/match"

// ModifyPowers is a middleware HandlerFunc to change the power of any card, not just itself
func ModifyPowers(h func(*match.GetPowerEvent)) func(*match.Card, *match.Context) {

	return func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.GetPowerEvent); ok {

			h(event)

		}

	}

}
