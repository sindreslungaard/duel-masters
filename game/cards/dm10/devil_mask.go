package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// NightmareInvader ...
func NightmareInvader(c *match.Card) {

	c.Name = "Nightmare Invader"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.DevilMask}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature)

}
