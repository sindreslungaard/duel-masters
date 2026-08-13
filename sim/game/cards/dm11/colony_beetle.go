package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SkyscraperShell ...
func SkyscraperShell(c *match.Card) {

	c.Name = "Skyscraper Shell"
	c.Power = 2000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.ColonyBeetle}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned, fx.OpponentChoosesCreatureToMana)))

}
