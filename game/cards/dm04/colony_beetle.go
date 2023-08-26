package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CannonShell ...
func CannonShell(c *match.Card) {

	c.Name = "Cannon Shell"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.PowerModifier = func(m *match.Match, attacking bool) int {

		shields, err := c.Player.Container(match.SHIELDZONE)

		if err != nil {
			return 0
		}

		return 1000 * len(shields)
	}

	c.Use(fx.Creature, fx.ShieldTrigger)

}
