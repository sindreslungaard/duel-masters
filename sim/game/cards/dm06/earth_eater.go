package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func HazardCrawler(c *match.Card) {

	c.Name = "Hazard Crawler"
	c.Power = 6000
	c.Civ = civ.Water
	c.Family = []string{family.EarthEater}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackCreatures, fx.CantAttackPlayers)
}

func MidnightCrawler(c *match.Card) {

	c.Name = "Midnight Crawler"
	c.Power = 6000
	c.Civ = civ.Water
	c.Family = []string{family.EarthEater}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Doublebreaker, fx.When(fx.Summoned, fx.ReturnOpCardFromMZToHand))
}

func ThrashCrawler(c *match.Card) {

	c.Name = "Thrash Crawler"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.EarthEater}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackCreatures, fx.CantAttackPlayers, fx.When(fx.Summoned, fx.ReturnMyCardFromMZToHand))
}
