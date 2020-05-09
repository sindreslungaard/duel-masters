package dm01

import (
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AquaHulcus ...
func AquaHulcus(c *match.Card) {

	c.Name = "Aqua Hulcus"
	// ...

	// Use existing middlewares, or define new ones
	c.Use(fx.Playable, fx.Creature, func(card *match.Card, c *match.Context) {

	})

}
