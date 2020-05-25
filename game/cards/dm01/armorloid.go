package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ArmoredWalkerUrherion ...
func ArmoredWalkerUrherion(c *match.Card) {

	c.Name = "Armored Walker Urherion"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = family.Armorloid
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

	c.PowerModifier = func(m *match.Match) int {

		power := 0

		if match.ContainerHas(c.Player, match.BATTLEZONE, func(x *match.Card) bool { return x.Family == family.Human }) {
			power += 2000
		}

		return power
	}

}
