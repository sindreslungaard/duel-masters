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
		fmt.Sprintf("%s's effect: Select 1 of your opponent's creatures and tap it.", card.Name),
		1,
		1,
		optional,
	).Map(func(creature *match.Card) {
		reportTapped(creature, ctx, card)
	})
}

// reportTapped taps a creature and says so in the chat. Tapping is otherwise
// silent, which leaves the affected player guessing at what just happened to
// their board.
func reportTapped(creature *match.Card, ctx *match.Context, source *match.Card) {
	creature.Tapped = true

	ctx.Match.ReportActionInChat(creature.Player, fmt.Sprintf("%s was tapped by %s", creature.Name, source.Name))
}

func TapOpCreature(card *match.Card, ctx *match.Context) {
	tapOpCreatureWithOptin(card, ctx, false)
}

func MayTapOpCreature(card *match.Card, ctx *match.Context) {
	tapOpCreatureWithOptin(card, ctx, true)
}

// TapUpToXOpCreatures taps up to x of the opponent's creatures, chosen by the
// card's controller. "Up to" makes the whole choice declinable, and a battle
// zone holding fewer than x creatures offers only the ones that are there.
func TapUpToXOpCreatures(x int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Choose up to %d of your opponent's creatures and tap them.", card.Name, x),
			1,
			x,
			true,
		).Map(func(creature *match.Card) {
			reportTapped(creature, ctx, card)
		})
	}
}
