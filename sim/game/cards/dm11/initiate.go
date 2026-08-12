package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// NialVizierOfDexterity ...
func NialVizierOfDexterity(c *match.Card) {

	c.Name = "Nial, Vizier of Dexterity"
	c.Power = 2500
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Initiate}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.EndOfMyTurnCreatureBZ, fx.MayUntapSelf))

}
