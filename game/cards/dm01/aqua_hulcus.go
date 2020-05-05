package dm01

import (
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AquaHulcus ...
func AquaHulcus(c *match.Card) {

	c.Name = "bob"

	c.Use(fx.Creature, func(card *match.Card, c *match.Context) {

	})

}
