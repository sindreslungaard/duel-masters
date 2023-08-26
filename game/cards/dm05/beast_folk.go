package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CrowWinger ...
func CrowWinger(c *match.Card) {

	c.Name = "Crow Winger"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		cards := fx.FindFilter(m.Opponent(c.Player), match.BATTLEZONE, func(x *match.Card) bool { return x.Civ == civ.Water || x.Civ == civ.Darkness })

		return len(cards) * 1000

	}

	c.Use(fx.Creature)

}
