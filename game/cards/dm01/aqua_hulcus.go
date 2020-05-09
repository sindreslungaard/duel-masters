package dm01

import (
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AquaHulcus ...
func AquaHulcus(c *match.Card) {

	c.Name = "Aqua Hulcus"
	c.Civ = "undefind_civ"
	c.Family = "undefined_family"
	c.ManaCost = 1
	c.ManaRequirement = []string{"test", "test2"}

	// Use existing middlewares, or define new ones
	c.Use(fx.Playable, fx.Creature)

}
