package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TulkTheOracle ...
func TulkTheOracle(c *match.Card) {

	c.Name = "Tulk, the Oracle"
	c.Civs = []string{civ.Light}
	c.Power = 500
	c.Family = []string{family.LightBringer}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature)

}

// KaemiraTheOracle ...
func KaemiraTheOracle(c *match.Card) {

	c.Name = "Kaemira, the Oracle"
	c.Power = 1000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.LightBringer}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.SilentSkill(fx.TopCardToShield))

}
