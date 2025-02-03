package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Gigagrax ...
func Gigagrax(c *match.Card) {

	c.Name = "Gigagrax"
	c.Power = 5000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Destroyed, fx.DestroyOpponentCreature(true, match.DestroyedByMiscAbility)))

}
