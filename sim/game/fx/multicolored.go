package fx

import (
	"duel-masters/game/match"
)

// PutIntoManaZoneTapped taps the card as it arrives in its owner's mana zone.
//
// This is the rule every multicolored card carries as reminder text, and it
// applies however the card gets there: charging mana, or an effect that puts it
// into the mana zone. A move always untaps the card first and only then
// dispatches CardMoved, so this handler has the last word on the tap state.
func PutIntoManaZoneTapped(card *match.Card, ctx *match.Context) {

	event, ok := ctx.Event.(*match.CardMoved)

	if !ok || event.CardID != card.ID || event.To != match.MANAZONE {
		return
	}

	card.Tapped = true

}
