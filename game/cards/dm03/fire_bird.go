package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BabyZoppe ...
func BabyZoppe(c *match.Card) {

	c.Name = "Baby Zoppe"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.FireBird}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		if match.ContainerHas(c.Player, match.MANAZONE, func(x *match.Card) bool { return x.Civ != civ.Fire }) {
			return 0
		}

		return 2000
	}

	c.Use(fx.Creature)
}
