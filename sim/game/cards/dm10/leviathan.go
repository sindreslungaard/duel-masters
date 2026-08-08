package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KingOquanos ...
func KingOquanos(c *match.Card) {
	c.Name = "King Oquanos"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		return len(fx.FindFilter(
			m.Opponent(c.Player),
			match.MANAZONE,
			func(card *match.Card) bool { return card.Tapped },
		)) * 2000
	}

	// King Oquanos only has the double breaker tier.
	c.Use(fx.Creature, fx.PowerBreakerTiers(6000, 0))
}
