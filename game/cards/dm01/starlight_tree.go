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
	c.Family = []string{family.StarlightTree}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)

}

// RubyGrass ...
func RubyGrass(c *match.Card) {

	c.Name = "Ruby Grass"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.StarlightTree}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers, fx.Untap)

}

// SenatineJadeTree ...
func SenatineJadeTree(c *match.Card) {

	c.Name = "Senatine Jade Tree"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.StarlightTree}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)

}
