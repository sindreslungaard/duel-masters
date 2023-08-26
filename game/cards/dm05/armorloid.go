package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CannoneerBargon ...
func CannoneerBargon(c *match.Card) {

	c.Name = "Cannoneer Bargon"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.CantAttackPlayers)

}
