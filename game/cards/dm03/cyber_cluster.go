package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AnglerCluster ...
func AnglerCluster(c *match.Card) {

	c.Name = "Angler Cluster"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.CyberCluster}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		if match.ContainerHas(c.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Water }) {
			return 0
		}

		return 3000
	}


	c.Use(fx.Creature, fx.Blocker, fx.CantAttackCreatures, fx.CantAttackPlayers)

}
