package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// LahPurificationEnforcer ...
func LahPurificationEnforcer(c *match.Card) {

	c.Name = "Lah, Purification Enforcer"
	c.Power = 5500
	c.Civ = civ.Light
	c.Family = family.Berserker
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}
