package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

func powerAttacker(card *match.Card, ctx *match.Context, n int) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		if ctx.Match.IsPlayerTurn(card.Player) {
			card.AddCondition(cnd.PowerAttacker, n, card.ID)
		}

	}

}

// PowerAttacker1000 adds the power attacker +1000 condition to the card
func PowerAttacker1000(card *match.Card, ctx *match.Context) {
	powerAttacker(card, ctx, 1000)
}

// PowerAttacker2000 adds the power attacker +2000 condition to the card
func PowerAttacker2000(card *match.Card, ctx *match.Context) {
	powerAttacker(card, ctx, 2000)
}

// PowerAttacker3000 adds the power attacker +3000 condition to the card
func PowerAttacker3000(card *match.Card, ctx *match.Context) {
	powerAttacker(card, ctx, 3000)
}

// PowerAttacker4000 adds the power attacker +4000 condition to the card
func PowerAttacker4000(card *match.Card, ctx *match.Context) {
	powerAttacker(card, ctx, 4000)
}

// PowerAttacker8000 adds the power attacker +8000 condition to the card
func PowerAttacker8000(card *match.Card, ctx *match.Context) {
	powerAttacker(card, ctx, 8000)
}
