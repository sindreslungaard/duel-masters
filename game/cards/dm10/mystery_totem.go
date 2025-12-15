package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// JigglyTotem ...
func JigglyTotem(c *match.Card) {

	c.Name = "Jiggly Totem"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.MysteryTotem}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if attacking {
			return len(fx.FindFilter(
				c.Player,
				match.MANAZONE,
				func(x *match.Card) bool {
					return x.Tapped
				},
			)) * 1000
		}

		return 0
	}

	c.Use(fx.Creature)

}
