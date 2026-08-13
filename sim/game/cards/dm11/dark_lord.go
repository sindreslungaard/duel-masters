package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HazariaDukeOfThorns ...
func HazariaDukeOfThorns(c *match.Card) {

	c.Name = "Hazaria, Duke of Thorns"
	c.Power = 2000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.DarkLord}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned, fx.OpponentChoosesAndDestroysCreature)))

}
