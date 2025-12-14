package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AdventureBoar ...
func AdventureBoar(c *match.Card) {

	c.Name = "Adventure Boar"
	c.Civ = civ.Nature
	c.Power = 1000
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}
