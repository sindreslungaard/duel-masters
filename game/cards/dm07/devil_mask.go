package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func ThreeFacedAshuraFang(c *match.Card) {

	c.Name = "Three-Faced Ashura Fang"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.DevilMask}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.PutShieldIntoHand))
}
