package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// LagunaLightningEnforcer ...
func LagunaLightningEnforcer(c *match.Card) {

	c.Name = "Laguna, Lightning Enforcer"
	c.Power = 2500
	c.Civ = civ.Light
	c.Family = []string{family.Berserker}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Attacking, fx.SearchDeckTake1Spell))

}
