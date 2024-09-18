package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KingCoral ...
func KingCoral(c *match.Card) {

	c.Name = "King Coral"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker)

}

// KingDepthcon ...
func KingDepthcon(c *match.Card) {

	c.Name = "King Depthcon"
	c.Power = 6000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.CantBeBlocked, fx.Creature, fx.Doublebreaker)

}

// KingRippedHide ...
func KingRippedHide(c *match.Card) {

	c.Name = "King Ripped-Hide"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.DrawUpTo2))

}
