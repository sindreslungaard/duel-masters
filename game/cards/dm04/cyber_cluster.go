package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HunterCluster ...
func HunterCluster(c *match.Card) {

	c.Name = "Hunter Cluster"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}


	c.Use(fx.Creature, fx.Blocker, fx.ShieldTrigger)

}
