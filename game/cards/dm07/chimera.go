package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func PhantasmalHorrorGigazabal(c *match.Card) {

	c.Name = "Phantasmal Horror Gigazabal"
	c.Power = 9000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.LightStealth)
}
