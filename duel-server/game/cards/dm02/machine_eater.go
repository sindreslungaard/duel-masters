package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// EngineerKipo ...
func EngineerKipo(c *match.Card) {

	c.Name = "Engineer Kipo"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.MachineEater}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.Destroyed, fx.EachPlayerDestroys1Mana))

}
