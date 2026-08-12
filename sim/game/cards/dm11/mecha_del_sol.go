package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// YulianaChannelerOfSuns ...
func YulianaChannelerOfSuns(c *match.Card) {

	c.Name = "Yuliana, Channeler of Suns"
	c.Power = 3000
	c.Civs = []string{civ.Light}
	c.Family = []string{family.MechaDelSol}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	// Being unchoosable does not make it untouchable: it can still be attacked
	// and blocked, which is exactly what CantBeSelectedByOpp models.
	c.Use(fx.Creature, fx.Blocker(), fx.CantBeSelectedByOpp, fx.CantAttackPlayers)

}
