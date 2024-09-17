package fx

import (
	"duel-masters/game/match"
	"fmt"
)

func draw(card *match.Card, ctx *match.Context, n int) {

	if event, ok := ctx.Event.(*match.CardMoved); ok {

		if event.CardID == card.ID && (event.To == match.BATTLEZONE || event.To == match.SPELLZONE) {

			card.Player.DrawCards(n)

		}

	}

}

// Draw1 draws 1 card when the card is added to the battlezone or spellzone
func Draw1(card *match.Card, ctx *match.Context) {
	draw(card, ctx, 1)
}

// Draw2 draws 2 card when the card is added to the battlezone or spellzone
func Draw2(card *match.Card, ctx *match.Context) {
	draw(card, ctx, 2)
}

// Draw3 draws 3 card when the card is added to the battlezone or spellzone
func Draw3(card *match.Card, ctx *match.Context) {
	draw(card, ctx, 3)
}

// Draw4 draws 4 card when the card is added to the battlezone or spellzone
func Draw4(card *match.Card, ctx *match.Context) {
	draw(card, ctx, 4)
}

// Draw5 draws 5 card when the card is added to the battlezone or spellzone
func Draw5(card *match.Card, ctx *match.Context) {
	draw(card, ctx, 5)
}

// Draw1ToMana draws 1 card and puts it in the players manazone
func Draw1ToMana(card *match.Card, ctx *match.Context) {

	cards := card.Player.PeekDeck(1)

	if len(cards) < 1 {
		return
	}

	c, err := card.Player.MoveCard(cards[0].ID, match.DECK, match.MANAZONE, card.ID)

	if err != nil {
		return
	}

	ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was added to %s's manazone from the top of their deck", c.Name, card.Player.Username()))

}

func Draw2ToMana(card *match.Card, ctx *match.Context) {
	Draw1ToMana(card, ctx)
	Draw1ToMana(card, ctx)
}

func MayDraw1(card *match.Card, ctx *match.Context) {
	MayDrawAmount(card, ctx, 1)
}

func DrawUpTo2(card *match.Card, ctx *match.Context) {
	DrawBetween(card, ctx, 0, 2)
}

func DrawUpTo3(card *match.Card, ctx *match.Context) {
	DrawBetween(card, ctx, 0, 3)
}

// This gives the player the choice to select a number of cards to draw between 2 provided limits
func DrawBetween(card *match.Card, ctx *match.Context, min int, max int) {
	count := max
	if min != max {
		count = SelectCount(
			card.Player,
			ctx.Match,
			fmt.Sprintf("%s effect: Draw between %d and %d cards", card.Name, min, max),
			min,
			max)
	}
	card.Player.DrawCards(count)
}

// This lets the player choose if they want to draw the full amount or none
func MayDrawAmount(card *match.Card, ctx *match.Context, amount int) {
	drawAmount := 0
	textAmount := fmt.Sprintf("%d cards", amount)
	if amount == 1 {
		textAmount = "1 card"
	}

	if BinaryQuestion(card.Player, ctx.Match, fmt.Sprintf("Do you want to draw %s? (%s effect)", textAmount, card.Name)) {
		drawAmount = amount
	}

	card.Player.DrawCards(drawAmount)
}
