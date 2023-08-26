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
	c.Family = []string{family.Hedrian}
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
	c.Family = []string{family.Hedrian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {

			if card.Zone != match.BATTLEZONE {

				exit()
				return

			}

			if event, ok := ctx2.Event.(*match.CardMoved); ok &&
				event.From == match.BATTLEZONE &&
				event.To == match.GRAVEYARD &&
				event.CardID != card.ID {

				if !ctx2.Match.IsPlayerTurn(card.Player) {
					ctx2.Match.Wait(ctx2.Match.Opponent(card.Player), "Waiting for your opponent to make an action")
					defer ctx2.Match.EndWait(ctx2.Match.Opponent(card.Player))
				}

				result := fx.SelectBacksideFilter(card.Player, ctx2.Match, card.Player, match.BATTLEZONE, fmt.Sprintf("%s: You may draw a card. Click close to not draw a card.", card.Name), 1, 1, true, func(x *match.Card) bool {
					return x.ID == c.ID
				})

				if len(result) > 0 {
					card.Player.DrawCards(1)
					ctx2.Match.Chat("Server", fmt.Sprintf("%s chose to draw a card from %s's ability", card.Player.Username(), card.Name))
				}

			}

		})

	}))

}
