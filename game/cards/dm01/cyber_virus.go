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
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.CantBeBlocked, fx.Creature)

}

// FaerieChild ...
func FaerieChild(c *match.Card) {

	c.Name = "Faerie Child"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.CantBeBlocked, fx.Creature)

}

// MarineFlower ...
func MarineFlower(c *match.Card) {

	c.Name = "Marine Flower"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, fx.CantAttackCreatures)

}
