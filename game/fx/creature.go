package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, ctx *match.Context) {

	// Resolve summoning sickness
	if _, ok := ctx.Event.(*match.BeginTurnStep); ok {

		if ctx.Match.IsPlayerTurn(card.Player) && card.HasCondition(cnd.SummoningSickness) {
			card.RemoveCondition(cnd.SummoningSickness)
		}

	}

	// Untap the card
	if _, ok := ctx.Event.(*match.UntapStep); ok {

		if ctx.Match.IsPlayerTurn(card.Player) {
			card.Tapped = false
		}

	}

	// Add to battlezone
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

				card.AddCondition(cnd.SummoningSickness)

				card.Player.MoveCard(card.ID, match.HAND, match.BATTLEZONE)

				break

			}

		})

	}

	// Attack the player
	if event, ok := ctx.Event.(*match.AttackPlayer); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		if card.HasCondition(cnd.SummoningSickness) {
			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot be played this turn as it has summoning sickness", card.Name))
			ctx.InterruptFlow()
			return
		}

		// Do this last in case any other cards want to interrupt the flow
		ctx.ScheduleAfter(func() {

			opponent := ctx.Match.Opponent(card.Player)

			// Allow the opponent to block if they can
			if len(event.Blockers) > 0 {

				ctx.Match.NewAction(opponent, event.Blockers, 1, 1, "You are being attacked. Choose a creature to block the attack with or close to not block the attack.", true)

				for {

					action := <-opponent.Action

					if action.Cancel {
						ctx.Match.CloseAction(opponent)
						break
					}

					if len(action.Cards) != 1 || !match.AssertCardsIn(event.Blockers, action.Cards[0]) {
						ctx.Match.ActionWarning(opponent, "Your selection of cards does not fulfill the requirements")
						continue
					}

					card, err := opponent.GetCard(action.Cards[0], match.BATTLEZONE)

					if err != nil {
						ctx.Match.ActionWarning(opponent, "The card you selected is not in the battlefield")
						continue
					}

					// ...

				}

			}

			shieldzone, err := opponent.Container(match.SHIELDZONE)

			if err != nil {
				return
			}

			if len(shieldzone) < 1 {
				// TODO: Attack the player, WIN if no blockers
			}

		})

	}

}
