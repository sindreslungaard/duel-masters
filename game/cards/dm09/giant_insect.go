package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// WhipScorpion ...
func WhipScorpion(c *match.Card) {

	c.Name = "Whip Scorpion"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.PowerAttacker3000)

}
