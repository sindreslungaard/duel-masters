package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ZeppelinCrawler ...
func ZeppelinCrawler(c *match.Card) {

	c.Name = "Zeppelin Crawler"
	c.Power = 4000
	c.Civ = civ.Water
	c.Family = []string{family.EarthEater}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackCreatures, fx.CantAttackPlayers,
		fx.LookTop4Put1IntoHandReorderRestOnBottomDeck)

}
