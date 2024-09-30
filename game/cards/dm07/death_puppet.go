package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func GezaryUndercoverDoll(c *match.Card) {

	c.Name = "Gezary, Undercover Doll"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.DeathPuppet}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.NatureStealth)
}
