package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CrystalLancer ...
func CrystalLancer(c *match.Card) {

	c.Name = "Crystal Lancer"
	c.Power = 8000
	c.Civ = civ.Water
	c.Family = family.LiquidPeople
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.CantBeBlocked, fx.Creature, fx.Evolution, fx.Doublebreaker)

}
