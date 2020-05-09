package fx

import (
	"duel-masters/game/match"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Playable has default functionality for all playable cards such as creatures and spells
func Playable(card *match.Card, c *match.Context) {

	// Check for and tap required mana when played, move to the battlezone
	if event, ok := c.Event.(*match.PlayCardEvent); ok {

		// Is this event for me or someone else?
		if event.CardID != card.ID {
			return
		}

		// OK, but do my player have the necessary amount of untapped mana?
		untappedMana := make([]*match.Card, 0)
		for _, card := range card.Player.Manazone {
			if !card.Tapped {
				untappedMana = append(untappedMana, card)
			}
		}

		// OK, but does the untapped mana contain at least 1 of each of the card's required civs?
		/* fulfillsManaRequirements := false
		for _, civ := range card.ManaRequirement {
			for _, card := range untappedMana {

			}
		} */

		c.Match.NewAction(
			card.Player,
			untappedMana,
			card.ManaCost,
			card.ManaCost,
			fmt.Sprintf("Select %v cards from your manazone to play %v. You must select at least 1 %v, civilization card.", card.ManaCost, card.Name, card.Civ),
			false,
		)

		for {

			action := <-card.Player.Action

			if ok := match.AssertCardsIn(untappedMana, action.Cards...); !ok {
				c.Match.WarnPlayer(card.Player, "You cannot select this card as it was not part of your given options")
				continue
			}

			logrus.Debug(action)

			break

		}

	}

}
