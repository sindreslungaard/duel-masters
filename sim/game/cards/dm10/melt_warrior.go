package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CoreCrashLizard ...
func CoreCrashLizard(c *match.Card) {

	c.Name = "Core-Crash Lizard"
	c.Power = 6000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.MeltWarrior}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	// The shield goes straight to the graveyard, so it is never broken and its
	// shield trigger is not offered.
	c.Use(fx.Creature, fx.When(fx.Summoned, fx.PutOpShieldIntoGraveyard))

}
