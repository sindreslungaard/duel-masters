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
	c.Family = family.Guardian
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackPlayers)

}
