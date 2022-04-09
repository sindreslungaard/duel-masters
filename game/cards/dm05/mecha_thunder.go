package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// LaGuileSeekerOfSkyfire ...
func LaGuileSeekerOfSkyfire(c *match.Card) {

	c.Name = "La Guile, Seeker of Skyfire"
	c.Power = 7500
	c.Civ = civ.Light
	c.Family = family.MechaThunder
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Doublebreaker)

}
