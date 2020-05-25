package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DeadlyFighterBraidClaw ...
func DeadlyFighterBraidClaw(c *match.Card) {

	c.Name = "Deadly Fighter Braid Claw"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = family.Dragonoid
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ForceAttack)

}
