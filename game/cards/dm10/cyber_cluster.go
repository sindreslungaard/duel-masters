package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BatteryCluster ...
func BatteryCluster(c *match.Card) {

	c.Name = "Battery Cluster"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackCreatures, fx.CantAttackPlayers)

}

// TorpedoCluster ...
func TorpedoCluster(c *match.Card) {

	c.Name = "Torpedo Cluster"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.ReturnMyCardFromMZToHand))

}
