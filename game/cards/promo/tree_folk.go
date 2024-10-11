package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// AngryMaple ...
func AngryMaple(c *match.Card) {

	c.Name = "Angry Maple"
	c.Power = 1000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.PowerAttacker4000)
}
