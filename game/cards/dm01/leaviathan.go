package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KingCoral ...
func KingCoral(c *match.Card) {

	c.Name = "King Coral"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = family.Leviathan
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker)

}
