package fx

import (
	"duel-masters/game/match"
	"fmt"
)

func Charger(card *match.Card, ctx *match.Context) {

	// Charger
	// After you cast this spell, put it into your mana zone instead of your graveyard.
	if event, ok := ctx.Event.(*match.SpellCast); ok && event.CardID == card.ID {

		card.Player.MoveCard(card.ID, match.HAND, match.MANAZONE)
		card.Tapped = false
		ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to the mana zone", card.Name))

	}
}
