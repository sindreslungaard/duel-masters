package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// FreiVizierOfAir ...
func FreiVizierOfAir(c *match.Card) {

	c.Name = "Frei, Vizier of Air"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = family.Initiate
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Untap)

}

// IereVizierOfBullets ...
func IereVizierOfBullets(c *match.Card) {

	c.Name = "Iere, Vizier of Bullets"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = family.Initiate
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}

// LokVizierOfHunting ...
func LokVizierOfHunting(c *match.Card) {

	c.Name = "Lok, Vizier of Hunting"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = family.Initiate
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}