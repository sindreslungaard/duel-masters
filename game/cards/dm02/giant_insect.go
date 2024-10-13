package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// XenoMantis ...
func XenoMantis(c *match.Card) {

	c.Name = "Xeno Mantis"
	c.Power = 6000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Attacking, fx.CantBeBlockedByPowerUpTo5000))

}
