package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func Aeropica(c *match.Card) {

	c.Name = "Aeropica"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}
	c.TapAbility = fx.ReturnCreatureToOwnersHand

	c.Use(fx.Creature, fx.TapAbility)
}

func Zepimeteus(c *match.Card) {

	c.Name = "Zepimeteus"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker, fx.CantAttackCreatures, fx.CantAttackPlayers)
}

func PromephiusQ(c *match.Card) {

	c.Name = "Promephius Q"
	c.Power = 2000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker, family.Survivor}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Survivor)
}
