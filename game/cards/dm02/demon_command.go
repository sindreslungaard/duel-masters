package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
	"math/rand"
)

// DarkTitanMaginn ...
func DarkTitanMaginn(c *match.Card) {

	c.Name = "Dark Titan Maginn"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {
		hand := fx.Find(ctx.Match.Opponent(card.Player), match.HAND)

		if len(hand) < 1 {
			return
		}

		discardedCard, err := ctx.Match.Opponent(card.Player).MoveCard(hand[rand.Intn(len(hand))].ID, match.HAND, match.GRAVEYARD)
		if err == nil {
			ctx.Match.Chat("Server", fmt.Sprintf("%s was discarded from %s's hand by Dark Titan Maginn", discardedCard.Name, discardedCard.Player.Username()))
		}
	}))

}

func dtmSpecial(card *match.Card, ctx *match.Context, cardID string) {

	if cardID != card.ID {
		return
	}

	ctx.ScheduleAfter(func() {

		hand, err := ctx.Match.Opponent(card.Player).Container(match.HAND)

		if err != nil {
			return
		}

		if len(hand) < 1 {
			return
		}

		discardedCard, err := ctx.Match.Opponent(card.Player).MoveCard(hand[rand.Intn(len(hand))].ID, match.HAND, match.GRAVEYARD)
		if err == nil {
			ctx.Match.Chat("Server", fmt.Sprintf("%s was discarded from %s's hand by %s", discardedCard.Name, discardedCard.Player.Username(), card.Name))
		}

	})

}
