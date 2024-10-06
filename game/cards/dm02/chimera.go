package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Gigastand ...
func Gigastand(c *match.Card) {

	c.Name = "Gigastand"
	c.Power = 3000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.WouldBeDestroyed, fx.MayReturnToHandAndDiscardACard))

}
