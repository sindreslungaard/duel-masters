package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func BroodShell(c *match.Card) {

	c.Name = "Brood Shell"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = fx.ReturnCreatureFromManazoneToHand

	c.Use(fx.Creature, fx.TapAbility)
}
