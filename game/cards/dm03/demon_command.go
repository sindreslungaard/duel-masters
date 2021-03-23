package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// GirielGhastlyWarrior ...
func GirielGhastlyWarrior(c *match.Card) {

	c.Name = "Giriel, Ghastly Warrior"
	c.Power = 11000
	c.Civ = civ.Darkness
	c.Family = family.DemonCommand
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Doublebreaker)

}
