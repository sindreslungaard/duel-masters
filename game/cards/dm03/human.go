package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// ExplosiveDudeJoe ...
func ExplosiveDudeJoe(c *match.Card) {

	c.Name = "Explosive Dude Joe"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

}

// MuramasaDukeOfBlades ...
func MuramasaDukeOfBlades(c *match.Card) {

	c.Name = "Muramasa, Duke of Blades"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: You may select 1 of your opponent's creatures with power 2000 or less and destroy it", card.Name),
			1,
			1,
			true,
			func(x *match.Card) bool { return ctx.Match.GetPower(x, false) <= 2000 },
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
		})
	}))

}
