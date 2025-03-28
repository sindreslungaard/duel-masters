package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// KyrstronLairDelver ...
func KyrstronLairDelver(c *match.Card) {

	c.Name = "Kyrstron, Lair Delver"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.ArmoredDragon}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.WouldBeDestroyed, fx.MayPutDragonFromHandIntoBZ))
}
