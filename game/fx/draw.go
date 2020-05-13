package fx

import "duel-masters/game/match"

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
