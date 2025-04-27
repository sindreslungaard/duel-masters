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

func MayReturnToHandAndDiscardACard(card *match.Card, ctx *match.Context) {
	if BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("%s was destroyed. Do you want to return it to hand? You will discard a card after", card.Name)) {
		ReturnToHand(card, ctx)
		Select(
			card.Player,
			ctx.Match,
			card.Player,
			match.HAND,
			"Select a card to discard",
			1,
			1,
			false,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.HAND, match.GRAVEYARD, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was discarded from %s's hand", x.Name, x.Player.Username()))
		})
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

func returnCreatureToOwnersHandWithOptin(card *match.Card, ctx *match.Context, optional bool) {
	cards := make(map[string][]*match.Card)

	cards["Your creatures"] = Find(card.Player, match.BATTLEZONE)
	cards["Opponent's creatures"] = Find(ctx.Match.Opponent(card.Player), match.BATTLEZONE)

	start_statement := "Choose"
	if optional {
		start_statement = "You may chooose"
	}

	SelectMultipart(
		card.Player,
		ctx.Match,
		cards,
		fmt.Sprintf("%s: %s 1 creature in the battlezone that will be sent to its owner's hand", card.Name, start_statement),
		1,
		1,
		optional,
	).Map(func(creature *match.Card) {
		creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.HAND, card.ID)
		ctx.Match.ReportActionInChat(creature.Player, fmt.Sprintf("%s was returned to %s's hand by %s", creature.Name, creature.Player.Username(), card.Name))
	})
}

func ReturnCreatureToOwnersHand(card *match.Card, ctx *match.Context) {
	returnCreatureToOwnersHandWithOptin(card, ctx, false)
}

func MayReturnCreatureToOwnersHand(card *match.Card, ctx *match.Context) {
	returnCreatureToOwnersHandWithOptin(card, ctx, true)
}

func PutOwnCreatureFromBZToMZ(card *match.Card, ctx *match.Context) {
	Select(card.Player, ctx.Match, card.Player, match.BATTLEZONE,
		"Select 1 of your creatures and put it in your manazone", 1, 1, false,
	).Map(func(creature *match.Card) {
		creature.Player.MoveCard(creature.ID, match.BATTLEZONE, match.MANAZONE, card.ID)
		creature.Tapped = false
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to manazone", creature.Name))
	})
}

func ReturnXCreaturesFromGraveToHand(x int) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.GRAVEYARD,
			fmt.Sprintf("%s: Return up to %d creature(s) from your graveyard to your hand", card.Name, x),
			1,
			x,
			false,
			func(x *match.Card) bool { return x.HasCondition(cnd.Creature) },
			true,
		).Map(func(x *match.Card) {
			card.Player.MoveCard(x.ID, match.GRAVEYARD, match.HAND, card.ID)
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was moved to %s's hand from their graveyard by %s", x.Name, card.Player.Username(), card.Name))
		})
	}
}
