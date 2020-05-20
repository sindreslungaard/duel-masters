package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AstrocometDragon ...
func AstrocometDragon(c *match.Card) {

	c.Name = "Astrocomet Dragon"
	c.Power = 6000
	c.Civ = civ.Fire
	c.Family = family.ArmoredDragon
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker4000, fx.Doublebreaker)

}
