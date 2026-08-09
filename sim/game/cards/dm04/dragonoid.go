package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BlastoExplosiveSoldier ...
func BlastoExplosiveSoldier(c *match.Card) {

	c.Name = "Blasto, Explosive Soldier"
	c.Power = 2000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		power := 0

		if match.ContainerHas(c.Player, match.BATTLEZONE, func(x *match.Card) bool { return x.HasCiv(civ.Darkness) }) {
			power += 2000
		}

		return power
	}

	c.Use(fx.Creature)
}
