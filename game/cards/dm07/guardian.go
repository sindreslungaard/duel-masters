package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func RodiGaleNightGuardian(c *match.Card) {

	c.Name = "Rodi Gale, Night Guardian"
	c.Power = 3500
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.DarknessStealth)
}
