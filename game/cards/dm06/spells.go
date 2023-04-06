package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func ProtectiveForce(c *match.Card) {

	c.Name = "Protective Force"
	c.Civ = civ.Light
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Spell, fx.When(fx.SpellCast, func(card *match.Card, ctx *match.Context) {

		fx.SelectFilter(
			card.Player,
			ctx.Match,
			card.Player,
			match.BATTLEZONE,
			"Protective Force: Select a blocker to give +4000 power to until the end of the turn",
			1,
			1,
			false,
			func(x *match.Card) bool {
				return x.HasCondition(cnd.Blocker)
			},
		).Map(func(x *match.Card) {
			x.AddCondition(cnd.PowerAmplifier, 4000, card.ID)
		})

		// TODO: make sure the condition is removed at the end of opponent's turn if shield trigger

	}))
}
