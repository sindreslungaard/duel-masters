package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, ctx *match.Context) {

	// Untap the card at the UntapStep
	if _, ok := ctx.Event.(*match.UntapStep); ok {

		if ctx.Match.IsPlayerTurn(card.Player) {
			card.Tapped = false
		}

	}

	// Check for and tap required mana when played, move to the battlezone
	if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		// Do this last in case any other cards want to interrupt the flow
		ctx.ScheduleAfter(func() {

			manazone, err := card.Player.Container(match.MANAZONE)

			if err != nil {
				return
			}

			untappedMana := make([]*match.Card, 0)
			for _, card := range manazone {
				if !card.Tapped {
					untappedMana = append(untappedMana, card)
				}
			}

			if !card.Player.CanPlayCard(card, untappedMana) {
				ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("You do not have sufficient mana to play %s", card.Name))
				ctx.InterruptFlow()
				return
			}

			ctx.Match.NewAction(
				card.Player,
				untappedMana,
				card.ManaCost,
				card.ManaCost,
				fmt.Sprintf("Select %v cards from your manazone to play %v. You must select at least 1 %v, civilization card.", card.ManaCost, card.Name, card.Civ),
				true,
			)

			defer ctx.Match.CloseAction(card.Player)

			for {

				action := <-card.Player.Action

				if action.Cancel {
					ctx.InterruptFlow()
					break
				}

				if len(action.Cards) != card.ManaCost || !match.AssertCardsIn(untappedMana, action.Cards...) {
					ctx.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
					continue
				}

				for _, id := range action.Cards {
					mana, err := card.Player.GetCard(id, match.MANAZONE)

					if err != nil {
						continue
					}

					mana.Tapped = true
				}

				card.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE)

				ctx.Match.BroadcastState()

				break

			}

		})

	}

}
