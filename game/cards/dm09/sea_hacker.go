package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Tekorax ...
func Tekorax(c *match.Card) {

	c.Name = "Tekorax"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.When(fx.Summoned, fx.LookAtOppShields))

}
