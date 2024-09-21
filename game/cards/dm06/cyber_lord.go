package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func Sopian(c *match.Card) {

	c.Name = "Sopian"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = fx.GiveOwnCreatureCantBeBlocked

	c.Use(fx.Creature, fx.TapAbility)
}
