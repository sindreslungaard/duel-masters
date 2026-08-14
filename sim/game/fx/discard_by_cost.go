package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// maxPrintedManaCost is the highest cost any card in the game has, and so the
// widest sensible range for "choose a number".
const maxPrintedManaCost = 13

// ChooseANumberAndDiscardByCost implements "Choose a number. Show your hand to
// your opponent and discard from it each card that has that cost. Then your
// opponent shows you his hand and discards from it each card that has that
// cost."
//
// The caster's own hand is emptied of that cost first, which matters because
// the number is chosen before either hand is seen.
func ChooseANumberAndDiscardByCost(card *match.Card, ctx *match.Context) {
	opponent := ctx.Match.Opponent(card.Player)

	cost := SelectCount(
		card.Player,
		ctx.Match,
		fmt.Sprintf("%s's effect: Choose a number. Both players discard every card of that cost from their hand.", card.Name),
		1,
		maxPrintedManaCost,
	)

	revealHandTo(ctx, card, card.Player, opponent)
	discardHandByCost(ctx, card, card.Player, cost)

	revealHandTo(ctx, card, opponent, card.Player)
	discardHandByCost(ctx, card, opponent, cost)
}

// revealHandTo shows one player's hand to another.
func revealHandTo(ctx *match.Context, source *match.Card, owner *match.Player, viewer *match.Player) {
	hand := Find(owner, match.HAND)

	if len(hand) < 1 {
		return
	}

	imageIDs := make([]string, 0, len(hand))
	for _, card := range hand {
		imageIDs = append(imageIDs, card.ImageID)
	}

	ctx.Match.ShowCards(viewer, fmt.Sprintf("%s's effect: %s's hand", source.Name, owner.Username()), imageIDs)
}

// discardHandByCost discards every card of the given cost from a player's hand.
func discardHandByCost(ctx *match.Context, source *match.Card, owner *match.Player, cost int) {
	// Snapshotted before anything moves, so the walk is not disturbed by the
	// hand shrinking underneath it.
	matching := FindFilter(owner, match.HAND, func(x *match.Card) bool { return x.ManaCost == cost })

	discarded := 0
	for _, card := range matching {
		moved, err := owner.MoveCard(card.ID, match.HAND, match.GRAVEYARD, source.ID)

		if err != nil || moved.Zone != match.GRAVEYARD {
			continue
		}

		discarded++
	}

	if discarded > 0 {
		ctx.Match.ReportActionInChat(owner, fmt.Sprintf("%s discarded %d card(s) that cost %d because of %s", owner.Username(), discarded, cost, source.Name))
	}
}
