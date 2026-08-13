package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// RevealTopXTake1ReorderRestOnBottom reveals the top x cards of the player's
// deck to both players, lets them take one that matches the filter into their
// hand, and puts everything else on the bottom of the deck in an order they
// choose.
//
// Revealing is not looking: the opponent sees the cards too. Nothing matching
// the filter is not a failure, it simply means the whole reveal goes to the
// bottom.
func RevealTopXTake1ReorderRestOnBottom(x int, filter func(*match.Card) bool, filterDescription string) match.HandlerFunc {
	return func(card *match.Card, ctx *match.Context) {
		revealed := card.Player.PeekDeck(x)

		if len(revealed) < 1 {
			return
		}

		imageIDs := make([]string, 0, len(revealed))
		for _, revealedCard := range revealed {
			imageIDs = append(imageIDs, revealedCard.ImageID)
		}

		message := fmt.Sprintf("%s's effect: the top %d cards of %s's deck", card.Name, len(revealed), card.Player.Username())
		ctx.Match.ShowCards(card.Player, message, imageIDs)
		ctx.Match.ShowCards(ctx.Match.Opponent(card.Player), message, imageIDs)

		taken := ""

		SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.DECK,
			fmt.Sprintf("%s's effect: Put one %s into your hand.", card.Name, filterDescription),
			1,
			1,
			false,
			func(x *match.Card) bool {
				if !filter(x) {
					return false
				}

				for _, revealedCard := range revealed {
					if revealedCard.ID == x.ID {
						return true
					}
				}

				return false
			},
			false,
		).Map(func(chosen *match.Card) {
			moved, err := card.Player.MoveCard(chosen.ID, match.DECK, match.HAND, card.ID)

			if err != nil || moved.Zone != match.HAND {
				return
			}

			taken = chosen.ID
			ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was put into %s's hand by %s", chosen.Name, card.Player.Username(), card.Name))
		})

		// Everything still in the deck goes to the bottom, including all of it
		// when nothing matched.
		rest := make([]*match.Card, 0, len(revealed))
		for _, revealedCard := range revealed {
			if revealedCard.ID != taken {
				rest = append(rest, revealedCard)
			}
		}

		if len(rest) < 1 {
			return
		}

		orderedIDs := OrderCards(
			card.Player,
			ctx.Match,
			rest,
			fmt.Sprintf("%s's effect: Order these cards that will be put on the bottom of your deck.", card.Name),
		)

		if len(orderedIDs) != len(rest) {
			return
		}

		card.Player.ReorderCardsInDeck(rest, orderedIDs, true)
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%d card(s) were put on the bottom of %s's deck by %s", len(rest), card.Player.Username(), card.Name))
	}
}

// ReturnToHandIfItBrokeShields implements "At the end of each of your turns, if
// this creature broke any shields that turn, return it to your hand."
//
// Whether the creature broke a shield is not recorded anywhere on the card, so
// the handler remembers it between the break and the end of the turn. The flag
// lives on the closure rather than in a condition because conditions are wiped
// at the end of every turn, which is exactly when this needs to read it.
func ReturnToHandIfItBrokeShields() match.HandlerFunc {
	brokeShields := false

	return func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.BrokenShieldEvent); ok {
			if event.Source == card.ID && card.Zone == match.BATTLEZONE {
				brokeShields = true
			}

			return
		}

		if !EndOfMyTurn(card, ctx) {
			return
		}

		// Reset first: whatever happens below, the next turn starts clean, and
		// a creature that left the battle zone in the meantime has nothing to
		// return.
		shouldReturn := brokeShields && card.Zone == match.BATTLEZONE
		brokeShields = false

		if !shouldReturn {
			return
		}

		moved, err := card.Player.MoveCard(card.ID, match.BATTLEZONE, match.HAND, card.ID)

		if err != nil || moved.Zone != match.HAND {
			return
		}

		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s returned to %s's hand after breaking shields this turn", card.Name, card.Player.Username()))
	}
}
