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

// DrawToMana draws 1 card and puts it in the players manazone
func DrawToMana(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.CardMoved); ok {

		if event.CardID == card.ID && (event.To == match.BATTLEZONE || event.To == match.SPELLZONE) {

			cards := card.Player.PeekDeck(1)

			if len(cards) < 1 {
				return
			}

			c, err := card.Player.MoveCard(cards[0].ID, match.DECK, match.MANAZONE, card.ID)

			if err != nil {
				return
			}

			ctx.Match.Chat("Server", fmt.Sprintf("%s was added to %s's manazone from the top of their deck", c.Name, ctx.Match.PlayerRef(card.Player).Socket.User.Username))

		}

	}

}

func MayDraw1(card *match.Card, ctx *match.Context) {

	ctx.Match.NewAction(card.Player, nil, 0, 0, "Do you want to draw a card?", true)

	action := <-card.Player.Action

	if !action.Cancel {
		draw(card, ctx, 1)
	}
	ctx.Match.CloseAction(card.Player)

}
