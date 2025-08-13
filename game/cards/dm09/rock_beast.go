package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Quakesaur ...
func Quakesaur(c *match.Card) {

	c.Name = "Quakesaur"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.WheneverThisAttacksPlayerAndIsntBlocked, fx.OpponentChoosesManaBurn))

}
