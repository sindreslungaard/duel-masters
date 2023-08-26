package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func WispHowlerShadowOfTears(c *match.Card) {
	c.Name = "Wisp Howler, Shadow of Tears"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.Ghost}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.ConditionalSlayer(func(target *match.Card) bool {
		return target.Civ == civ.Nature || target.Civ == civ.Light
	}))
}
