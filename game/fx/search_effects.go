package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

func SearchDeckMoveCardsZone(card *match.Card, ctx *match.Context, quantity int, filter func(*match.Card) bool, playerText string, newZone string, resultText string) {

	SelectFilter(card.Player,
		ctx.Match,
		card.Player,
		match.DECK,
		fmt.Sprintf("%s effect: %s", card.Name, playerText),
		1,
		quantity,
		true,
		filter,
		true,
	).Map(func(x *match.Card) {
		x.Player.MoveCard(x.ID, match.DECK, newZone, card.ID)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf(resultText, x.Name))
	})

	ShuffleDeck(card, ctx, false)

}

func SearchDeckPutIntoManazone(card *match.Card, ctx *match.Context, quantity int, filter func(*match.Card) bool, filterDescription string) {

	SearchDeckMoveCardsZone(
		card,
		ctx,
		quantity,
		filter,
		fmt.Sprintf("You make move up to %d %s from your deck to your manazone", quantity, filterDescription),
		match.MANAZONE,
		card.Player.Username()+" put %s from their deck to their manazone",
	)

}

func SearchDeckTakeCards(card *match.Card, ctx *match.Context, quantity int, filter func(*match.Card) bool, filterDescription string) {
	SearchDeckMoveCardsZone(
		card,
		ctx,
		quantity,
		filter,
		fmt.Sprintf("You make take up to %d %s from your deck (Will be shown to your opponent)", quantity, filterDescription),
		match.HAND,
		card.Player.Username()+" retrieved %s from their deck to their hand",
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
