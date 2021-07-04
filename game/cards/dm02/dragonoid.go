package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CavalryGeneralCuratops ...
func CavalryGeneralCuratops(c *match.Card) {

	c.Name = "Cavalry General Curatops"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = family.Dragonoid
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.AttackUntapped)

}
