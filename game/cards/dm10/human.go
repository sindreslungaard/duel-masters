package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MezgerCommandoLeader ...
func MezgerCommandoLeader(c *match.Card) {

	c.Name = "Mezger Commando Leader"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Human}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.SpeedAttacker)

}
