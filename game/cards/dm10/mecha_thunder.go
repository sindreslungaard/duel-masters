package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BalzaSeekerOfHyperpearls ...
func BalzaSeekerOfHyperpearls(c *match.Card) {

	c.Name = "Balza, Seeker of Hyperpearls"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.MechaThunder}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.Blocker(), fx.CantAttackPlayers)

}
