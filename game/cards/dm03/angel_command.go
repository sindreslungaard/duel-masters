package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MiarCometElemental ...
func MiarCometElemental(c *match.Card) {

	c.Name = "Miar, Comet Elemental"
	c.Power = 11500
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Doublebreaker, fx.Creature)

}
