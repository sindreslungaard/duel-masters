package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HypersprintWariorUzesol ...
func HypersprintWariorUzesol(c *match.Card) {

	c.Name = "Hypersprint Warior Uzesol"
	c.Power = 1000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.SpeedAttacker, fx.PowerAttacker4000)

}
