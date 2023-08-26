package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// OuksVizierOfRestoration ...
func OuksVizierOfRestoration(c *match.Card) {

	c.Name = "Ouks, Vizier of Restoration"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.ReturnToShield)

}

// SariusVizierOfSuppression ...
func SariusVizierOfSuppression(c *match.Card) {

	c.Name = "Sarius, Vizier of Suppression"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)

}
