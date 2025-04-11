package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MagmadragonMelgars ...
func MagmadragonMelgars(c *match.Card) {

	c.Name = "Magmadragon Melgars"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.VolcanoDragon}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature)

}
