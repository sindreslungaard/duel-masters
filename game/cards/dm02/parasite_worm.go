package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ChaosWorm ...
func ChaosWorm(c *match.Card) {

	c.Name = "Chaos Worm"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.Select(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: Select a creature from the opponent's battle zone and destroy it", card.Name),
			1,
			1,
			true,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
		})
	}))

}

// UltracideWorm ...
func UltracideWorm(c *match.Card) {

	c.Name = "Ultracide Worm"
	c.Power = 11000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker)

}

// HorridWorm ...
func HorridWorm(c *match.Card) {

	c.Name = "Horrid Worm"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, fx.OpponentDiscardsRandomCard))

}

// PoisonWorm ...
func PoisonWorm(c *match.Card) {

	c.Name = "Poison Worm"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			fmt.Sprintf("%s: Destroy one of your creatures with power 3000 or less", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 3000 },
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
		})
	}))

}
