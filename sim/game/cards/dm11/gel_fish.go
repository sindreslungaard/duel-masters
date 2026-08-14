package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// FantasyFish ...
func FantasyFish(c *match.Card) {

	c.Name = "Fantasy Fish"
	c.Power = 2000
	c.Civs = []string{civ.Water}
	c.Family = []string{family.GelFish}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.ShieldTrigger, fx.Blocker())

}
