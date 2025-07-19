package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// StingerBall ...
func StingerBall(c *match.Card) {

	c.Name = "Stinger Ball"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.CyberVirus}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.ShieldsSelectionEffect, fx.WheneverThisAttacksMayLookAtOpShield())

}
