package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func OtherworldlyWarriorNaglu(c *match.Card) {

	c.Name = "Otherworldly Warrior Naglu"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.Armorloid}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker3000, fx.Doublebreaker, fx.CantBeAttacked)

}
