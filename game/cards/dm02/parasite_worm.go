package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
	"math/rand"
)

// ChaosWorm ...
func ChaosWorm(c *match.Card) {

	c.Name = "Chaos Worm"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = family.ParasiteWorm
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {

		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			"Chaos Worm: Select a creature from the opponent's battle zone and destroy it",
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card)
		})

	}))

}

// UltracideWorm ...
func UltracideWorm(c *match.Card) {

	c.Name = "Ultracide Worm"
	c.Power = 11000
	c.Civ = civ.Darkness
	c.Family = family.ParasiteWorm
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker)

}

// HorridWorm ...
func HorridWorm(c *match.Card) {

	c.Name = "Horrid Worm"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = family.ParasiteWorm
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Attacking, func(card *match.Card, ctx *match.Context) {

		hand := fx.Find(ctx.Match.Opponent(card.Player), match.HAND)

		if len(hand) < 1 {
			return
		}

		discardedCard, err := ctx.Match.Opponent(card.Player).MoveCard(hand[rand.Intn(len(hand))].ID, match.HAND, match.GRAVEYARD)
		if err == nil {
			ctx.Match.Chat("Server", fmt.Sprintf("%s was discarded from %s's hand by Horrid Worm", discardedCard.Name, discardedCard.Player.Username()))
		}

	}))

}
