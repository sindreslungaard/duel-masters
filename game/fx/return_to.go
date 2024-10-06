package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// ReturnToHand returns the card to the players hand instead of the graveyard
func ReturnToHand(card *match.Card, ctx *match.Context) {
	ctx.InterruptFlow()

	card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND, card.ID)
	ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was returned to the hand", card.Name))
}

func MayReturnToHand(card *match.Card, ctx *match.Context) {
	if BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("%s was destroyed. Do you want to return it to hand", card.Name)) {
		ReturnToHand(card, ctx)
	}
}

// ReturnToMana returns the card to the players manazone instead of the graveyard
func ReturnToMana(card *match.Card, ctx *match.Context) {

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

		if event.Card == card {

			ctx.InterruptFlow()

			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.MANAZONE, card.ID)
			card.Tapped = false
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was destroyed by %s and moved to the mana zone", event.Card.Name, event.Source.Name))

		}

	}

}

// ReturnToShield returns the card to the players shield zone instead of the graveyard
func ReturnToShield(card *match.Card, ctx *match.Context) {

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {

		if event.Card == card {

			ctx.InterruptFlow()

			card.Player.MoveCard(card.ID, match.BATTLEZONE, match.SHIELDZONE, card.ID)
			card.Tapped = false
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was destroyed by %s and moved to the shield zone", event.Card.Name, event.Source.Name))

		}

	}

}

// PutShieldIntoHand Player picks an own shield and puts it into their hand
func PutShieldIntoHand(card *match.Card, ctx *match.Context) {
	SelectBackside(
		card.Player,
		ctx.Match,
		card.Player,
		match.SHIELDZONE,
		fmt.Sprintf("%s: Move 1 of your shields into your hand.", card.Name),
		1,
		1,
		false,
	).Map(func(x *match.Card) {
		ctx.Match.MoveCard(x, match.HAND, card)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s effect: shield moved to hand", card.Name))
	})
}

func ReturnCreatureFromManazoneToHand(card *match.Card, ctx *match.Context) {
	SelectFilter(card.Player, ctx.Match, card.Player, match.MANAZONE,
		"Select 1 of your creatures from your mana zone that will be returned to your hand",
		1, 1, false, func(x *match.Card) bool { return x.HasCondition(cnd.Creature) }, false,
	).Map(func(x *match.Card) {
		card.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's hand from their mana zone", x.Name, ctx.Match.PlayerRef(card.Player).Socket.User.Username))
	})
}
