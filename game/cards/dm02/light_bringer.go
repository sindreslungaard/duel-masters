package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// FonchTheOracle ...
func FonchTheOracle(c *match.Card) {

	c.Name = "Fonch, the Oracle"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.TapOpCreature))

}

// WynTheOracle ...
func WynTheOracle(c *match.Card) {

	c.Name = "Wyn, the Oracle"
	c.Power = 1500
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.ShieldsSelectionEffect, fx.WheneverThisAttacksMayLookAtOpShield())

}
