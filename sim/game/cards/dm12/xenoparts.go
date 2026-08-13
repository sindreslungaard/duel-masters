package dm12

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MuramasasKnife ...
func MuramasasKnife(c *match.Card) {

	c.Name = "Muramasa's Knife"
	c.Power = 2000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Xenoparts}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.AttackUntapped)

}

// GandavalsStapler ...
func GandavalsStapler(c *match.Card) {

	c.Name = "Gandaval's Stapler"
	c.Power = 3000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Xenoparts}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.TapSelfWhenAnotherCreatureArrives)

}
