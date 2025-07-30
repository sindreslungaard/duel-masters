package fx

import (
	"duel-masters/game/cnd"
	"duel-masters/game/match"
	"fmt"
)

// Spell has default functionality for spells
func Spell(card *match.Card, ctx *match.Context) {

	// Clear conditions
	if _, ok := ctx.Event.(*match.EndOfTurnStep); ok {
		card.ClearConditions()
	}

	// Add spell condition to this card
	if _, ok := ctx.Event.(*match.UntapStep); ok {
		card.AddCondition(cnd.Spell, nil, card.ID)
	}

	// When the spell is played from hand
	if event, ok := ctx.Event.(*match.PlayCardEvent); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		// make sure we haven't attacked yet
		if _, ok := ctx.Match.Step.(*match.AttackStep); ok {
			ctx.Match.WarnPlayer(card.Player, "You can't cast spells after attacking")
			ctx.InterruptFlow()
			return
		}

		ctx.ScheduleAfter(func() {

			manazone, err := card.Player.Container(match.MANAZONE)

			if err != nil {
				return
			}

			untappedMana := make([]*match.Card, 0)
			for _, c := range manazone {
				if !c.Tapped {
					untappedMana = append(untappedMana, c)
				}
			}

			if !card.Player.CanPlayCard(card, untappedMana) {
				ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("You do not have sufficient mana to play %s", card.Name))
				ctx.InterruptFlow()
				return
			}

			manaCost := card.ManaCost
			for _, condition := range card.Conditions() {
				if condition.ID == cnd.ReducedCost {
					manaCost -= condition.Val.(int)
					if manaCost < 1 {
						manaCost = 1
					}
				}

				if condition.ID == cnd.IncreasedCost {
					manaCost += condition.Val.(int)
				}
			}

			ctx.Match.NewAction(
				card.Player,
				untappedMana,
				manaCost,
				manaCost,
				fmt.Sprintf("Select %v cards from your manazone to play %v. You must select at least 1 %v, civilization card.", manaCost, card.Name, card.Civ),
				true,
			)

			for {

				action := <-card.Player.Action

				if action.Cancel {
					ctx.Match.CloseAction(card.Player)
					ctx.InterruptFlow()
					break
				}

				cards := make([]*match.Card, 0)

				for _, id := range action.Cards {
					mana, err := card.Player.GetCard(id, match.MANAZONE)

					if err != nil {
						continue
					}

					cards = append(cards, mana)
				}

				if len(action.Cards) != manaCost || !match.AssertCardsIn(untappedMana, action.Cards...) || !card.Player.CanPlayCard(card, cards) {
					ctx.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
					continue
				}

				ctx.Match.CloseAction(card.Player)

				for _, mana := range cards {
					mana.Tapped = true
				}

				ctx.Match.CastSpell(card, false)

				break

			}

		})

	}

	// On spell cast
	if event, ok := ctx.Event.(*match.SpellCast); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s played the spell %s", card.Player.Username(), card.Name))

		ctx.ScheduleAfter(func() {
			card.Player.MoveCard(card.ID, match.HAND, match.GRAVEYARD, card.ID)
		})

	}

}
