package fx

import (
	"duel-masters/game/match"
	"fmt"
)

func Charger(card *match.Card, ctx *match.Context) {
	// After you cast this spell, put it into your mana zone instead of your
	// graveyard. fx.Spell already moved it to the graveyard by the time it
	// finished resolving, so relocate it from there.
	if event, ok := ctx.Event.(*match.SpellResolved); ok && event.CardID == card.ID {
		card.Player.MoveCard(card.ID, match.GRAVEYARD, match.MANAZONE, card.ID)
		card.Tapped = false
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was put in the mana zone instead of your graveyard", card.Name))
	}
}
