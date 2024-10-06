package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func StingerHornTheDelver(c *match.Card) {

	c.Name = "Stinger Horn, the Delver"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.HornedBeast}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker1000, fx.WaterStealth)
}
