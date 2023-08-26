package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GregorianWorm ...
func GregorianWorm(c *match.Card) {

	c.Name = "Gregorian Worm"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.ShieldTrigger)

}
