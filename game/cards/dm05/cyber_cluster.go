package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func SteelTurretCluster(c *match.Card) {
	c.Name = "Steel-Turret Cluster"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(
		fx.Creature,
		fx.CantBeAttackedIf(func(attacker *match.Card) bool {
			return attacker.Civ == civ.Fire || attacker.Civ == civ.Nature
		}),
	)
}
