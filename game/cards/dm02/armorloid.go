package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DogarnTheMarauder ...
func DogarnTheMarauder(c *match.Card) {

	c.Name = "Dogarn, the Marauder"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		power := 0

		if attacking {

			fx.FindFilter(
				c.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.Tapped && x != c },
			).Map(func(x *match.Card) {
				power += 2000
			})

		}

		return power

	}

	c.Use(fx.Creature)

}
