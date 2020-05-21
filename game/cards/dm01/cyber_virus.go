package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CandyDrop ...
func CandyDrop(c *match.Card) {

	c.Name = "Candy Drop"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = family.CyberVirus
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.CantBeBlocked, fx.Creature)

}
