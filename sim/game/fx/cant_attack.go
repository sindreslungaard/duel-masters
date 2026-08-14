package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// CantAttackPlayers prevents a card from attacking players
func CantAttackPlayers(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddCondition(cnd.CantAttackPlayers, true, card.ID)
	}

}

// CantAttackCreatures prevents a card from attacking players
func CantAttackCreatures(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddCondition(cnd.CantAttackCreatures, true, card.ID)
	}

}

func CantAttackIfOpHasNoShields(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {
		handleCantAttackIfOpHasNoShieldsConditions(card, ctx)
	}

	if event, ok := ctx.Event.(*match.CardMoved); ok {
		if event.From == match.SHIELDZONE || event.To == match.SHIELDZONE {
			handleCantAttackIfOpHasNoShieldsConditions(card, ctx)
		}
	}
}

func handleCantAttackIfOpHasNoShieldsConditions(card *match.Card, ctx *match.Context) {
	opponentShields := Find(ctx.Match.Opponent(card.Player), match.SHIELDZONE)
	n := len(opponentShields)

	if n == 0 {
		card.AddCondition(cnd.CantAttackPlayers, true, card.ID)
		card.AddCondition(cnd.CantAttackCreatures, true, card.ID)
	} else {
		card.RemoveCondition(cnd.CantAttackPlayers)
		card.RemoveCondition(cnd.CantAttackCreatures)
	}
}
