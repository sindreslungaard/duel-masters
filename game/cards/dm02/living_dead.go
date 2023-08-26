package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MarrowOozeTheTwister ...
func MarrowOozeTheTwister(c *match.Card) {

	c.Name = "Marrow Ooze, the Twister"
	c.Power = 1000
	c.Civ = civ.Darkness
	c.Family = []string{family.LivingDead}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker, fx.When(fx.AttackingPlayer, func(card *match.Card, ctx *match.Context) {

		ctx.ScheduleAfter(func() {
			ctx.Match.Destroy(card, card, match.DestroyedByMiscAbility)
		})

	}))

}
