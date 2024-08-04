package dm02

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// DarkTitanMaginn ...
func DarkTitanMaginn(c *match.Card) {

	c.Name = "Dark Titan Maginn"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.DemonCommand}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.AttackConfirmed, fx.OpponentDiscardsRandomCard))

}
