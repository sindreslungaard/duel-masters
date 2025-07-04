package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CandyCluster ...
func CandyCluster(c *match.Card) {

	c.Name = "Candy Cluster"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.CantBeBlocked)

}
