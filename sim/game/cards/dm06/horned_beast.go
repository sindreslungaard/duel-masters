package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func FeatherHornTheTracker(c *match.Card) {

	c.Name = "Feather Horn, the Tracker"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature)
}

func ParadiseHorn(c *match.Card) {

	c.Name = "Paradise Horn"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)
}
