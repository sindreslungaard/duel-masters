package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func Cratersaur(c *match.Card) {
	c.Name = "Cratersaur"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		condition1 := &match.Condition{ID: cnd.AttackUntapped, Val: true, Src: nil}
		condition2 := &match.Condition{ID: cnd.PowerAttacker, Val: 3000, Src: nil}
		fx.HaveSelfConditionsWhenNoShields(card, ctx, []*match.Condition{condition1, condition2})
	}))
}
