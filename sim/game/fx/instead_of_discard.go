package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// PutIntoBattleZoneInsteadOfDiscard implements "When this creature would be
// discarded from your hand during your opponent's turn, you may put it into the
// battle zone instead."
//
// Cancelling the pending move is what replaces the discard; the creature is
// then put into play out of the hand it never left. Only a move from hand to
// graveyard counts, so charging it as mana or playing it normally is untouched.
func PutIntoBattleZoneInsteadOfDiscard(card *match.Card, ctx *match.Context) {
	event, ok := ctx.Event.(*match.MoveCard)

	if !ok ||
		event.CardID != card.ID ||
		event.From != match.HAND ||
		event.To != match.GRAVEYARD ||
		ctx.Match.IsPlayerTurn(card.Player) {
		return
	}

	if !BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("Do you want to put %s into the battle zone instead of discarding it?", card.Name)) {
		return
	}

	ctx.InterruptFlow()

	ForcePutCreatureIntoBZ(ctx, card, match.HAND, card)
}

// TapSelfWhenAnotherCreatureArrives implements "Whenever another creature is
// put into the battle zone, tap this creature." Either player's creature counts.
func TapSelfWhenAnotherCreatureArrives(card *match.Card, ctx *match.Context) {
	event, ok := ctx.Event.(*match.CardMoved)

	if !ok || card.Zone != match.BATTLEZONE || event.To != match.BATTLEZONE || event.CardID == card.ID {
		return
	}

	// An evolution base resurfacing is not a creature arriving.
	if event.From == match.HIDDENZONE {
		return
	}

	card.Tapped = true
	ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was tapped because another creature entered the battle zone", card.Name))
}
