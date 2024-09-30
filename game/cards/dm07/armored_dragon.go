package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func GazariasDragon(c *match.Card) {
	c.Name = "Gazarias Dragon"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		condition1 := &match.Condition{ID: cnd.DoubleBreaker, Val: true, Src: nil}
		condition2 := &match.Condition{ID: cnd.PowerAmplifier, Val: 4000, Src: nil}
		fx.HaveSelfConditionsWhenNoShields(card, ctx, []*match.Condition{condition1, condition2})
	}))

}
