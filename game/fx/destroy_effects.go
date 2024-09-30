package fx

import (
	"duel-masters/game/match"
	"fmt"
)

func EachPlayerDestroys1Mana(card *match.Card, ctx *match.Context) {
	EachPlayerDestroysMana(card, ctx, 1)
}

func EachPlayerDestroysMana(card *match.Card, ctx *match.Context, quantity int) {

	players := make([]*match.Player, 0)
	players = append(players, card.Player)
	players = append(players, ctx.Match.Opponent(card.Player))

	for _, p := range players {

		cards := len(Find(p, match.MANAZONE))
		if quantity > cards {
			quantity = cards
		}

		Select(
			p,
			ctx.Match,
			p,
			match.MANAZONE,
			fmt.Sprintf("%s effect: Select %v card(s) from your manazone that will be sent to your graveyard", card.Name, quantity),
			quantity,
			quantity,
			false,
		).Map(func(manaCard *match.Card) {
			p.MoveCard(manaCard.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(p, fmt.Sprintf("%s effect: %s moved from MZ to GY", card.Name, manaCard.Name))
		})

	}

}
