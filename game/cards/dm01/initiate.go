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