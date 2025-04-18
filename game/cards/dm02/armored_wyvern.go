package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// MetalwingSkyterror ...
func MetalwingSkyterror(c *match.Card) {

	c.Name = "Metalwing Skyterror"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredWyvern}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {
		fx.SelectFilter(
			card.Player,
			ctx.Match,
			ctx.Match.Opponent(card.Player),
			match.BATTLEZONE,
			fmt.Sprintf("%s: You may destroy one of your opponent's blockers", card.Name),
			1,
			1,
			true,
			func(x *match.Card) bool { return x.HasCondition(cnd.Blocker) },
			false,
		).Map(func(x *match.Card) {
			ctx.Match.Destroy(x, card, match.DestroyedByMiscAbility)
		})
	}))

}
