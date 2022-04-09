package dm05

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BolgashDragon ...
func BolgashDragon(c *match.Card) {

	c.Name = "Bolgash Dragon"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = family.ArmoredDragon
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Triplebreaker, fx.PowerAttacker8000)

}
