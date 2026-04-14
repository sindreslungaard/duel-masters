package dm03

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BaragaBladeOfGloom ...
func BaragaBladeOfGloom(c *match.Card) {

	c.Name = "Baraga, Blade of Gloom"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.DarkLord}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Summoned, fx.PutShieldIntoHand))

}
