package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// Player chooses and mana burns x opponent's cards
func ManaBurnX(x int) func(*match.Card, *match.Context) {
	return func(card *match.Card, ctx *match.Context) {
		Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.MANAZONE,
			fmt.Sprintf("Select up to %d card(s) from your opponent's mana zone that will be sent to their graveyard", x),
			1,
			x,
			true,
		).Map(func(mana *match.Card) {
			mana.Player.MoveCard(mana.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was sent from %s's manazone to their graveyard by %s", mana.Name, mana.Player.Username(), card.Name))
		})
	}
}
