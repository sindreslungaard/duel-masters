package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BurningMane ...
func BurningMane(c *match.Card) {

	c.Name = "BurningMane"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = family.BeastFolk
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature)

}

// FearFang ...
func FearFang(c *match.Card) {

	c.Name = "Fear Fang"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = family.BeastFolk
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature)

}
