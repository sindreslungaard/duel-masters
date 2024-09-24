package fx

import (
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

func ReturnOpCardFromMZToHand(card *match.Card, ctx *match.Context) {
	Select(
		card.Player,
		ctx.Match,
		ctx.Match.Opponent(card.Player),
		match.MANAZONE,
		fmt.Sprintf("%s: Choose a card from your opponent's mana zone that will be returned to his hand.", card.Name),
		1,
		1,
		false,
	).Map(func(x *match.Card) {
		x.Player.MoveCard(x.ID, match.MANAZONE, match.HAND, card.ID)
		ctx.Match.ReportActionInChat(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s got moved to %s hand from his mana zone by %s", x.Name, x.Player.Username(), card.Name))
	})
}

func ReturnMyCardFromMZToHand(card *match.Card, ctx *match.Context) {
	Select(
		card.Player,
		ctx.Match,
		card.Player,
		match.MANAZONE,
		"Select 1 card from your mana zone that will be sent to your hand",
		1,
		1,
		false,
	).Map(func(c *match.Card) {
		c.Player.MoveCard(c.ID, match.MANAZONE, match.HAND, card.ID)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s effect: %s returned %s from MZ to their hand", card.Name, c.Player.Username(), c.Name))
	})
}
