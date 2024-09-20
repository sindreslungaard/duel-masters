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
	c.Power = 8000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker2000)

}

// StampedingLonghorn ...
func StampedingLonghorn(c *match.Card) {

	c.Name = "Stampeding Longhorn"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Attacking, fx.CantBeBlockedByPowerUpTo3000))

}

// TrihornShepherd ...
func TrihornShepherd(c *match.Card) {

	c.Name = "Tri-horn Shepherd"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature)

}
