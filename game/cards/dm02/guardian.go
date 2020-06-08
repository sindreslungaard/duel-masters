package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// LadiaBaleTheInspirational ...
func LadiaBaleTheInspirational(c *match.Card) {

	c.Name = "Ladia Bale, the Inspirational"
	c.Power = 9500
	c.Civ = civ.Light
	c.Family = family.Guardian
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker)

}
