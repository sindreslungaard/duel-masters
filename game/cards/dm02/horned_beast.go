package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RumblingTerahorn ...
func RumblingTerahorn(c *match.Card) {

	c.Name = "Rumbling Terahorn"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.SearchDeckTake1Creature))

}

// LeapingTornadoHorn ...
func LeapingTornadoHorn(c *match.Card) {

	c.Name = "Leaping Tornado Horn"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
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
