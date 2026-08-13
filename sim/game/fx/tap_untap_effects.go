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

func tapOpCreatureWithOptin(card *match.Card, ctx *match.Context, optional bool) {
	Select(
		card.Player,
		ctx.Match,
		ctx.Match.Opponent(card.Player),
		match.BATTLEZONE,
		"Select 1 of your opponent's creature and tap it.",
		1,
		1,
		optional,
	).Map(func(creature *match.Card) {
		creature.Tapped = true
	})
}

func TapOpCreature(card *match.Card, ctx *match.Context) {
	tapOpCreatureWithOptin(card, ctx, false)
}

func MayTapOpCreature(card *match.Card, ctx *match.Context) {
	tapOpCreatureWithOptin(card, ctx, true)
}
