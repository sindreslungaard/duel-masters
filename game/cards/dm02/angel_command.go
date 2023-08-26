package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// EthelStarSeaElemental ...
func EthelStarSeaElemental(c *match.Card) {

	c.Name = "Ethel, Star Sea Elemental"
	c.Power = 5500
	c.Civ = civ.Light
	c.Family = []string{family.AngelCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.CantBeBlocked, fx.Creature)

}
