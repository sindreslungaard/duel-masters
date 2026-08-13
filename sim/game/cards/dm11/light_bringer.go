package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MerleeTheOracle ...
func MerleeTheOracle(c *match.Card) {

	c.Name = "Merlee, the Oracle"
	c.Power = 1500
	c.Civs = []string{civ.Light}
	c.Family = []string{family.LightBringer}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	// "Each of your creatures" includes Merlee itself, so the grant is not
	// filtered down to the others.
	c.Use(fx.Creature, fx.WaveStriker, fx.WaveStrikerGrantToOwnCreatures(cnd.PowerAmplifier, 1000))

}
