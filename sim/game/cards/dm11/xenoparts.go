package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BradsCutter ...
func BradsCutter(c *match.Card) {

	c.Name = "Brad's Cutter"
	c.Power = 1000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Xenoparts}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ShieldTrigger)

}

// JabahasAutomaton ...
func JabahasAutomaton(c *match.Card) {

	c.Name = "Jabaha's Automaton"
	c.Power = 6000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Xenoparts}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Evolution, fx.PowerAttacker4000, fx.Doublebreaker)

}
