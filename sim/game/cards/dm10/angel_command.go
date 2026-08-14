package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ElixiaPurebladeElemental ...
func ElixiaPurebladeElemental(c *match.Card) {
	c.Name = "Elixia, Pureblade Elemental"
	c.Power = 1000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.PowerModifier = func(m *match.Match, attacking bool) int {
		if c.Zone != match.BATTLEZONE {
			return 0
		}

		// Each civilization counts once no matter how many cards of it are in
		// the mana zone. A multicolored card contributes every civilization on it.
		civilizations := make(map[string]struct{})
		for _, manaCard := range fx.Find(c.Player, match.MANAZONE) {
			for _, civilization := range manaCard.Civs {
				civilizations[civilization] = struct{}{}
			}
		}

		return len(civilizations) * 3000
	}

	c.Use(fx.Creature, fx.PowerBreakerTiers(6000, 15000))
}
