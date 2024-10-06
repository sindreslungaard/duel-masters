package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func LaunchLocust(c *match.Card) {

	c.Name = "Launch Locust"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}
	c.PowerModifier = func(m *match.Match, attacking bool) int {

		power := 0

		if attacking {
			for _, creature := range fx.Find(c.Player, match.BATTLEZONE) {
				if creature == c {
					continue
				}
				power += 1000
			}
		}

		return power
	}

	c.Use(fx.Creature)
}
