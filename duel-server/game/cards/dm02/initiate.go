package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MagrisVizierOfMagnetism ...
func MagrisVizierOfMagnetism(c *match.Card) {

	c.Name = "Magris, Vizier of Magnetism"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.MayDraw1))

}
