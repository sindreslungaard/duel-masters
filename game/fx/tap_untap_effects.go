package fx

import (
	"duel-masters/game/match"
	"fmt"
)

func MayUntapSelf(card *match.Card, ctx *match.Context) {

	if !card.Tapped {
		return
	}
	if BinaryQuestion(
		card.Player,
		ctx.Match,
		fmt.Sprintf("%s effect: Do you want to untap self?", card.Name),
	) {
		card.Tapped = false
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s untapped self", card.Name))
	}
}
