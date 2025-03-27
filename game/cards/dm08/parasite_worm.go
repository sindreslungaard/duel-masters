package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TyrantWorm ...
func TyrantWorm(c *match.Card) {

	c.Name = "Tyrant Worm"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.AnotherOwnCreatureSummoned, fx.DestroyYourself))
}
