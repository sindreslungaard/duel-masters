package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func KingTriumphant(c *match.Card) {

	c.Name = "King Triumphant"
	c.Power = 7000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Doublebreaker, ReceiveBlockerWhenOpponentPlaysCreatureOrSpell)
}
