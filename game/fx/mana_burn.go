package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// Player chooses and mana burns x opponent's cards
func ManaBurnX(x int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.MANAZONE,
<<<<<<< HEAD
			fmt.Sprintf("Select up to %d card(s) from your opponent's mana zone that will be sent to their graveyard", x),
=======
			fmt.Sprintf("Select upto %d card(s) from your opponent's mana zone that will be sent to their graveyard", x),
>>>>>>> Implement Dragon Evolution. Add Bajula and Abzo Dolba.
			1,
			x,
			false,
		).Map(func(mana *match.Card) {
			mana.Player.MoveCard(mana.ID, match.MANAZONE, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was sent from %s's manazone to their graveyard by %s", mana.Name, mana.Player.Username(), card.Name))
		})
	}
}
