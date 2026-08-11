package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TagtappTheRetaliator ...
func TagtappTheRetaliator(c *match.Card) {

	c.Name = "Tagtapp, the Retaliator"
	c.Power = 3000
	c.Civs = []string{civ.Fire, civ.Nature}
	c.Family = []string{family.SpiritQuartz}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire, civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		// A multicolored card in the opponent's mana zone counts as a water card
		// whenever water is one of its civilizations.
		return len(fx.FindFilter(
			m.Opponent(c.Player),
			match.MANAZONE,
			func(x *match.Card) bool { return x.HasCiv(civ.Water) },
		)) * 1000
	}

	c.Use(fx.Creature, fx.PowerBreakerTiers(6000, 0))

}
