package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// LothRix ...
func LothRix(c *match.Card) {

	c.Name = "Loth Rix, the Iridescent"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Evolution, fx.When(fx.Summoned, fx.TopCardToShield))

}
