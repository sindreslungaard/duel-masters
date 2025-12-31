package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RikabuTheDismantler ...
func RikabuTheDismantler(c *match.Card) {

	c.Name = "Rikabu, the Dismantler"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.MachineEater}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.SpeedAttacker)

}
