package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
	"math/rand"
)

// Locomotiver ...
func Locomotiver(c *match.Card) {

	c.Name = "Locomotiver"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = family.Hedrian
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

		if err != nil {
			return
		}

		if len(hand) > 0 {
			discardedCard, err := ctx.Match.Opponent(card.Player).MoveCard(hand[rand.Intn(len(hand))].ID, match.HAND, match.GRAVEYARD)
			if err == nil {
				ctx.Match.Chat("Server", fmt.Sprintf("%s was discarded from %s's hand", discardedCard.Name, discardedCard.Player.Username()))
			}
		}
	}))
}

// MongrelMan ...
func MongrelMan(c *match.Card) {

	c.Name = "Mongrel Man"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = family.Hedrian
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.CreatureDestroyed); ok {
			if event.Card.ID != card.ID {
				card.Player.DrawCards(1)
			}
		}

	})
}
