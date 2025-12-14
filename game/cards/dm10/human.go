package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ArmoredRaiderGandaval ...
func ArmoredRaiderGandaval(c *match.Card) {

	c.Name = "Armored Raider Gandaval"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if attacking {
			return 2000 * len(fx.FindFilter(
				c.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return x.ID != c.ID && x.Tapped
				},
			))
		}

		return 0
	}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker)

}
