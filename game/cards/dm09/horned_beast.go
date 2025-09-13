package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SolidHorn ...
func SolidHorn(c *match.Card) {

	c.Name = "Solid Horn"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.WouldBeDestroyed, fx.ReturnToMana))

}
