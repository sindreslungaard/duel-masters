package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CalgoVizierOfRainclouds ...
func CalgoVizierOfRainclouds(c *match.Card) {

	c.Name = "Calgo, Vizier of Rainclouds"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.CantBeBlockedByPower4000OrMore)

}
