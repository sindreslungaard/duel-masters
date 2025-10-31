package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// PoltalesterTheSpydroid ...
func PoltalesterTheSpydroid(c *match.Card) {

	c.Name = "Poltalester, the Spydroid"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Soltrooper}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.ShieldTrigger, fx.CantAttackPlayers)

}
