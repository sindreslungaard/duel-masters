package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ArtisanPicora ...
func ArtisanPicora(c *match.Card) {

	c.Name = "Artisan Picora"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = family.MachineEater
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.DestroyManaOnSummon)

}
