package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

func searchDeckMoveCardsZone(card *match.Card, ctx *match.Context, quantity int, filter func(*match.Card) bool, playerText string, newZone string, resultText string, showRetrievedCard bool) {
	SelectFilter(
		card.Player,
		ctx.Match,
		card.Player,
		match.DECK,
		fmt.Sprintf("%s effect: %s", card.Name, playerText),
		0,
		quantity,
		true,
		filter,
		true,
	).Map(func(x *match.Card) {
		card.Player.MoveCard(x.ID, match.DECK, newZone, card.ID)
		if showRetrievedCard {
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf(resultText, x.Name))
		} else {
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s retrieved a card from their deck to their hand", card.Player.Username()))
		}
	})

	ShuffleDeck(card, ctx, false)
}

func SearchDeckPutIntoManazone(card *match.Card, ctx *match.Context, quantity int, filter func(*match.Card) bool, filterDescription string) {
	searchDeckMoveCardsZone(
		card,
		ctx,
		quantity,
		filter,
		fmt.Sprintf("You may move up to %d %s from your deck to your manazone", quantity, filterDescription),
		match.MANAZONE,
		card.Player.Username()+" put %s from their deck to their manazone",
		true,
	)
}

func SearchDeckTakeCards(card *match.Card, ctx *match.Context, quantity int, filter func(*match.Card) bool, filterDescription string) {
	searchDeckMoveCardsZone(
		card,
		ctx,
		quantity,
		filter,
		fmt.Sprintf("You may take up to %d %s from your deck (Will be shown to your opponent)", quantity, filterDescription),
		match.HAND,
		card.Player.Username()+" retrieved %s from their deck to their hand",
		true,
	)
}

func SearchDeckTake1Spell(card *match.Card, ctx *match.Context) {
	SearchDeckTakeCards(
		card,
		ctx,
		1,
		func(x *match.Card) bool { return x.HasCondition(cnd.Spell) },
		"spell",
	)
}

func SearchDeckTake1Creature(card *match.Card, ctx *match.Context) {
	SearchDeckTakeCards(
		card,
		ctx,
		1,
		func(x *match.Card) bool { return x.HasCondition(cnd.Creature) },
		"creature",
	)
}

// Search deck for x cards
func SearchDeckTakeXCardsWithoutShowing(x int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		searchDeckMoveCardsZone(
			card,
			ctx,
			x,
			func(c *match.Card) bool { return true },
			fmt.Sprintf("You make take up to %d cards from your deck", x),
			match.HAND,
			"",
			false,
		)
	}
}
