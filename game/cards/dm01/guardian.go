package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DiaNorkMoonlightGuardian ...
func DiaNorkMoonlightGuardian(c *match.Card) {

	c.Name = "Dia Nork, Moonlight Guardian"
	c.Power = 5000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)

}

// GranGureSpaceGuardian ...
func GranGureSpaceGuardian(c *match.Card) {

	c.Name = "Gran Gure, Space Guardian"
	c.Power = 9000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)

}

// LaUraGigaSkyGuardian ...
func LaUraGigaSkyGuardian(c *match.Card) {

	c.Name = "La Ura Giga, Sky Guardian"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)

}

// SzubsKinTwilightGuardian ...
func SzubsKinTwilightGuardian(c *match.Card) {

	c.Name = "Szubs Kin, Twilight Guardian"
	c.Power = 6000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)

}
