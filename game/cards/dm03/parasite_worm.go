package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// HangWormFetidLarva ...
func HangWormFetidLarva(c *match.Card) {

	c.Name = "Hang Worm, FetidLarva"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = family.ParasiteWorm
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature)

}
