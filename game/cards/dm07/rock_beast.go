package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BlazosaurQ ...
func ValkrowzerUltraRockBeast(c *match.Card) {

	c.Name = "Valkrowzer, Ultra Rock Beast"
	c.Power = 9000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Evolution, fx.WaterStealth, fx.Doublebreaker)

}
