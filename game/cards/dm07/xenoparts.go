package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func KiposContraption(c *match.Card) {

	c.Name = "Kipo's Contraption"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.Xenoparts}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}
	c.TapAbility = fx.DestroyByMiscOpCreature2000OrLess

	c.Use(fx.Creature, fx.TapAbility)
}
