package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// PeppiPepper ...
func PeppiPepper(c *match.Card) {

	c.Name = "Peppi Pepper"
	c.Power = 2000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.FireBird}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker3000)

}

// BuzzBetocchi ...
func BuzzBetocchi(c *match.Card) {

	c.Name = "Buzz Betocchi"
	c.Power = 4000
	c.Civs = []string{civ.Fire, civ.Nature}
	c.Family = []string{family.FireBird, family.GiantInsect}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire, civ.Nature}

	c.Use(fx.Creature, fx.PutIntoManaZoneTapped)

}
