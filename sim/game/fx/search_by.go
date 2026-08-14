package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// SearchDeckPutCreatureIntoBZ lets the player look through their deck for a
// creature matching the filter and put it into the battle zone. The search is
// optional, and the deck is shuffled afterwards either way.
func SearchDeckPutCreatureIntoBZ(card *match.Card, ctx *match.Context, filter func(*match.Card) bool, description string) {
	SelectFilter(
		card.Player,
		ctx.Match,
		card.Player,
		match.DECK,
		fmt.Sprintf("%s's effect: You may put %s from your deck into the battle zone.", card.Name, description),
		1,
		1,
		true,
		func(x *match.Card) bool { return x.HasCondition(cnd.Creature) && filter(x) },
		true,
	).Map(func(found *match.Card) {
		ForcePutCreatureIntoBZ(ctx, found, match.DECK, card)
	})

	ShuffleDeck(card, ctx, false)
}
