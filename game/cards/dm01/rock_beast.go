package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Meteosaur ...
func Meteosaur(c *match.Card) {

	c.Name = "Meteosaur"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.DestroyByMiscOpCreature2000OrLess))

}

// Stonesaur ...
func Stonesaur(c *match.Card) {

	c.Name = "Stonesaur"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}
