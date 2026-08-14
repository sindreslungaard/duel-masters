package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// EngbeltTheSpydroid ...
func EngbeltTheSpydroid(c *match.Card) {

	c.Name = "Engbelt, the Spydroid"
	c.Power = 5500
	c.Civs = []string{civ.Light}
	c.Family = []string{family.Soltrooper}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers)

}
