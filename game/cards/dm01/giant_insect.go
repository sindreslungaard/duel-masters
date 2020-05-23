package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DeathbladeBeetle ...
func DeathbladeBeetle(c *match.Card) {

	c.Name = "Deathblade Beetle"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = family.GiantInsect
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Doublebreaker, fx.PowerAttacker4000)

}
