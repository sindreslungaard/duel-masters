package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// DestroyManaOnSummon forces the user to destroy one mana when the card is summoned
func DestroyManaOnSummon(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.CardMoved); ok {

		if event.CardID == card.ID && (event.To == match.BATTLEZONE || event.To == match.SPELLZONE) {

			manazone, err := card.Player.Container(match.MANAZONE)

			if err != nil {
				return
			}

			if len(manazone) < 1 {
				return
			}

			ctx.Match.NewAction(card.Player, manazone, 1, 1, "Select 1 card from your manazone that will be sent to your graveyard", false)

			for {

				action := <-card.Player.Action

				if len(action.Cards) != 1 || !match.AssertCardsIn(manazone, action.Cards[0]) {
					ctx.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
					continue
				}

				c, err := card.Player.MoveCard(action.Cards[0], match.MANAZONE, match.GRAVEYARD)

				if err != nil {
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's manazone to the graveyard", c.Name, ctx.Match.PlayerRef(card.Player).Socket.User.Username))
				}

				break

			}

		}

	}

}
