package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MishaChannelerOfSuns ...
func MishaChannelerOfSuns(c *match.Card) {

	c.Name = "Misha, Channeler of Suns"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.CantBeAttackedIf(func(attacker *match.Card) bool {
		return attacker.SharesAFamily(family.Dragons)
	}))
}
