package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AquaStrummer ...
func AquaStrummer(c *match.Card) {

	c.Name = "Aqua Strummer"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.LookAtUpTo5CardsFromTopDeckAndReorder))

}

// CrystalSpinslicer ...
func CrystalSpinslicer(c *match.Card) {

	c.Name = "Crystal Spinslicer"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Evolution, fx.Blocker())

}
