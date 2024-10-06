package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func TitaniumCluster(c *match.Card) {

	c.Name = "Titanium Cluster"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantBeAttacked, fx.CantAttackCreatures, fx.CantAttackPlayers)
}
