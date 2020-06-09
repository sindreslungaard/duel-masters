package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ScissorEye ...
func ScissorEye(c *match.Card) {

	c.Name = "Scissor Eye"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = family.GelFish
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature)

}
