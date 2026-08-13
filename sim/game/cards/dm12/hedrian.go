package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// WindmillMutant ...
func WindmillMutant(c *match.Card) {

	c.Name = "Windmill Mutant"
	c.Power = 2000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Hedrian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, fx.OpponentDiscardsRandomCard))

}

// SteamrollerMutant ...
func SteamrollerMutant(c *match.Card) {

	c.Name = "Steamroller Mutant"
	c.Power = 3000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.Hedrian}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	// "Destroy all creatures" spares nothing, Steamroller Mutant included, and
	// the wave strikers that switched its ability on go with it.
	c.Use(fx.Creature, fx.WaveStriker, fx.WhileWaveStriker(fx.When(fx.Summoned,
		fx.DestroyAllCreatures(match.DestroyedByMiscAbility))))

}
