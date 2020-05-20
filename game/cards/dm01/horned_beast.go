package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RoaringGreatHorn ...
func RoaringGreatHorn(c *match.Card) {

	c.Name = "Roaring Great-Horn"
	c.Power = 6000
	c.Civ = civ.Nature
	c.Family = family.BeastFolk
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker2000)

}
