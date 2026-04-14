package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func PicorasWrench(c *match.Card) {

	c.Name = "Picora's Wrench"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Xenoparts}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)
}

func RikabusScrewdriver(c *match.Card) {

	c.Name = "Rikabu's Screwdriver"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.Xenoparts}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}
	c.TapAbility = func(card *match.Card, ctx *match.Context) {
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
	}

	c.Use(fx.Creature, fx.TapAbility)
}
