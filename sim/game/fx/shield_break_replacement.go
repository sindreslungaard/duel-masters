package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// ShieldsToGraveyardInsteadOfHand implements "Whenever this creature would
// break a shield, your opponent puts that shield into his graveyard instead."
//
// Interrupting BreakShieldEvent stops the default move to hand entirely,
// which also means the shield never reaches hand and so never offers its
// shield trigger, matching the printed replacement.
//
// Register with fx.When(fx.BreakShield, fx.ShieldsToGraveyardInsteadOfHand),
// which already guards the event type and this card's battle zone presence.
func ShieldsToGraveyardInsteadOfHand(card *match.Card, ctx *match.Context) {

	event, ok := ctx.Event.(*match.BreakShieldEvent)

	if !ok || event.Source != card {
		return
	}

	ctx.InterruptFlow()

	for _, shield := range event.Cards {
		moved, err := ctx.Match.Opponent(card.Player).MoveCard(shield.ID, match.SHIELDZONE, match.GRAVEYARD, card.ID)
		if err == nil {
			ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s was moved to %s's graveyard instead of hand by %s", moved.Name, ctx.Match.Opponent(card.Player).Username(), card.Name))
		}
	}

}
