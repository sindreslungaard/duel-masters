package dm01

import (
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AquaHulcus ...
func AquaHulcus() *match.Card {

	c := &match.Card{}

	c.Use(fx.Creature)

	return c

}
