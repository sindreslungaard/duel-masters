package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// Cragsaur ...
func Cragsaur(c *match.Card) {

	c.Name = "Cragsaur"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

}

// Hurlosaur ...
func Hurlosaur(c *match.Card) {

	c.Name = "Hurlosaur"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.When(fx.Summoned, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s's effect: Destroy one of your opponent's creatures that was power 1000 or less.", card.Name),
			1,
			1,
			false,
			func(x *match.Card) bool {
				return ctx.Match.GetPower(x, false) <= 1000
			},
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
		})
	}))

}
