package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func LeQuistTheOracle(c *match.Card) {
	c.Name = "Le Quist, the Oracle"
	c.Power = 1500
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.WheneverThisAttacksMayTapDorFCreature())

}
