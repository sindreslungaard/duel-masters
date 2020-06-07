package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MiniTitanGett ...
func MiniTitanGett(c *match.Card) {

	c.Name = "Mini Titan Gett"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = family.Human
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ForceAttack, fx.PowerAttacker1000)

}
