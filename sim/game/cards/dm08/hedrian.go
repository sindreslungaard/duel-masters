package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MotorcycleMutant ...
func MotorcycleMutant(c *match.Card) {

	c.Name = "Motorcycle Mutant"
	c.Power = 6000
	c.Civ = civ.Darkness
	c.Family = []string{family.Hedrian}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackCreatures, fx.CantAttackPlayers,
		fx.When(fx.AnotherOwnCreatureSummoned, fx.DestroyYourself),
	)

}
