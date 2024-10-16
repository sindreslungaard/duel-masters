package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ArmoredGroblav	...
func ArmoredGroblav(c *match.Card) {

	c.Name = "Armored Groblav"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		power := 0

		if attacking {

			fx.FindFilter(
				c.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool { return x.Civ == civ.Fire && x != c },
			).Map(func(x *match.Card) {
				power += 1000
			})

		}

		return power

	}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker)

}
