package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SabermaskScarab ...
func SabermaskScarab(c *match.Card) {

	c.Name = "Sabermask Scarab"
	c.Power = 4000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.WheneverThisAttacksReturnCardFromMZToHand())

}
