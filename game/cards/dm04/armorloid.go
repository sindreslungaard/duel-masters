package dm04

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KamikazeChainsawWarrior ...
func KamikazeChainsawWarrior(c *match.Card) {

	c.Name = "Kamikaze, Chainsaw Warrior"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ShieldTrigger)
}
