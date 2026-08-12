package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// RoyalDurian ...
func RoyalDurian(c *match.Card) {

	c.Name = "Royal Durian"
	c.Power = 1000
	c.Civs = []string{civ.Nature}
	c.Family = []string{family.WildVeggies}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.SilentSkill(dragonFromMana()))

}
