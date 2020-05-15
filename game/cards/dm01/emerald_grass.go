package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// EmeraldGrass ...
func EmeraldGrass(c *match.Card) {

	c.Name = "Emerald Grass"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = family.StarlightTree
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)

}
