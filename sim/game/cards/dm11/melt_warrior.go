package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HysteriaLizard ...
func HysteriaLizard(c *match.Card) {

	c.Name = "Hysteria Lizard"
	c.Power = 3000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.MeltWarrior}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ForceAttack, fx.PowerAttacker3000)

}

// BonfireLizard ...
func BonfireLizard(c *match.Card) {

	c.Name = "Bonfire Lizard"
	c.Power = 4000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.MeltWarrior}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned, fx.DestroyUpToXOpBlockers(2))))

}
