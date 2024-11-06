package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func KizarBasikuTheOutrageous(c *match.Card) {

	c.Name = "Kizar Basiku, the Outrageous"
	c.Power = 8500
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker, fx.Evolution, fx.FireStealth, fx.Doublebreaker)
}

func RomVizierofTendrils(c *match.Card) {

	c.Name = "Rom, Vizier of Tendrils"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.MayTapOpCreature))
}
