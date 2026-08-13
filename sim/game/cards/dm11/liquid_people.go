package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AquaTrickster ...
func AquaTrickster(c *match.Card) {

	c.Name = "Aqua Trickster"
	c.Power = 1000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.LiquidPeople}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned, fx.TapOpCreature)))

}
