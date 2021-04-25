package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ExplosiveDudeJoe ...
func ExplosiveDudeJoe(c *match.Card) {

	c.Name = "Explosive Dude Joe"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = family.Human
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

}
