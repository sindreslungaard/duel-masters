package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TentacleCluster ...
func TentacleCluster(c *match.Card) {

	c.Name = "Tentacle Cluster"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.WheneverThisAttacksPlayerAndIsntBlocked, fx.ReturnCreatureToOwnersHand))

}
