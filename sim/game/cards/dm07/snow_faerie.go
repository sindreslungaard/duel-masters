package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func PoppleFlowerpetalDancer(c *match.Card) {

	c.Name = "Popple, Flowerpetal Dancer"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.SnowFaerie}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = fx.Draw1ToMana

	c.Use(fx.Creature, fx.TapAbility)
}
