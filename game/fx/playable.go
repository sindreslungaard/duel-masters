package fx

import (
	"duel-masters/game/match"
	"fmt"
)

// Playable has default functionality for all playable cards such as creatures and spells
func Playable(card *match.Card, c *match.Context) {

	// Check for and tap required mana when played, move to the battlezone
	if event, ok := c.Event.(*match.PlayCardEvent); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

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
			c.Match.WarnPlayer(card.Player, fmt.Sprintf("You do not have sufficient mana to play %s", card.Name))
			c.InterruptFlow()
			return
		}

		c.Match.NewAction(
			card.Player,
			untappedMana,
			card.ManaCost,
			card.ManaCost,
			fmt.Sprintf("Select %v cards from your manazone to play %v. You must select at least 1 %v, civilization card.", card.ManaCost, card.Name, card.Civ),
			true,
		)

		defer c.Match.CloseAction(card.Player)

		for {

			action := <-card.Player.Action

			if action.Cancel {
				c.InterruptFlow()
				break
			}

			if len(action.Cards) != card.ManaCost || !match.AssertCardsIn(untappedMana, action.Cards...) {
				c.Match.ActionWarning(card.Player, "Your selection of cards does not fulfill the requirements")
				continue
			}

			for _, id := range action.Cards {
				mana, err := card.Player.GetCard(id, match.MANAZONE)

				if err != nil {
					continue
				}

				mana.Tapped = true
			}

			c.Match.BroadcastState()

			break

		}

	}

}
