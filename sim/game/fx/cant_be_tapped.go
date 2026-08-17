package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
)

// CantBeTappedByOpp keeps the card out of any effect that lets the opponent
// choose it to tap, such as TapOpCreature or TapUpToXOpCreatures. It does not
// stop the card from being tapped by its own controller's effects, or as a
// cost or consequence of battling.
func CantBeTappedByOpp(card *match.Card, ctx *match.Context) {

	if _, ok := ctx.Event.(*match.UntapStep); ok {

		card.AddCondition(cnd.CantBeTappedByOpp, nil, card.ID)

	}

}
