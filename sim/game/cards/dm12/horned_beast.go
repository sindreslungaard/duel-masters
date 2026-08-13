package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RadioactiveHornTheStrange ...
func RadioactiveHornTheStrange(c *match.Card) {

	c.Name = "Radioactive Horn, the Strange"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker)

}

// SpectralHornGlitalis ...
func SpectralHornGlitalis(c *match.Card) {

	c.Name = "Spectral Horn Glitalis"
	c.Power = 4000
	c.Civs = []string{civ.Nature, civ.Light}
	c.Family = []string{family.HornedBeast, family.RainbowPhantom}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature, civ.Light}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped)

}
