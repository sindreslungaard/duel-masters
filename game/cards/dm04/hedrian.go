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

				if !ctx.Match.IsPlayerTurn(card.Player) {
					ctx.Match.Wait(ctx.Match.Opponent(card.Player), "Waiting for your opponent to make an action")
					defer ctx.Match.EndWait(ctx.Match.Opponent(card.Player))
				}

				result := fx.SelectBacksideFilter(card.Player, ctx.Match, card.Player, match.BATTLEZONE, "Mongrel Man: You may optionally draw a card because a creature was destroyed. Click close to not draw a card.", 1, 1, true, func(x *match.Card) bool {
					return x.ID == event.Card.ID
				})

				if len(result) > 0 {
					card.Player.DrawCards(1)
					ctx.Match.Chat("Server", fmt.Sprintf("%s chose to draw a card from Mongrel Man's ability", c.Player.Username()))
				}
			}

		}
	})
}
