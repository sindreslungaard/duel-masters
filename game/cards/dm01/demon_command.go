package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DeathligerLionOfChaos ...
func DeathligerLionOfChaos(c *match.Card) {

	c.Name = "Deathliger, Lion of Chaos"
	c.Power = 9000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker)

}

// ZagaanKnightOfDarkness ...
func ZagaanKnightOfDarkness(c *match.Card) {

	c.Name = "Zagaan, Knight of Darkness"
	c.Power = 7000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker)

}
