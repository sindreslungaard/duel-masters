package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

func PromephiusQ(c *match.Card) {

	c.Name = "Promephius Q"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker, family.Survivor}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature)

}
