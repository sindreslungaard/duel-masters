package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HazardHopper ...
func HazardHopper(c *match.Card) {

	c.Name = "Hazard Hopper"
	c.Power = 5000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.ReturnToHandIfItBrokeShields())

}
