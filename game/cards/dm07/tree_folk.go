package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func WorldTreeRootOfLife(c *match.Card) {

	c.Name = "World Tree, Root of Life"
	c.Power = 7000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution, fx.PowerAttacker2000, fx.DarknessStealth, fx.Doublebreaker)

}
