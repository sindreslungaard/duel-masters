package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RazaVegaThunderGuardian ...
func RazaVegaThunderGuardian(c *match.Card) {

	c.Name = "Raza Vega, Thunder Guardian"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 10
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.ReturnToShield)
}
