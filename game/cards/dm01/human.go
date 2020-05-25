package dm01

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// BrawlerZyler ...
func BrawlerZyler(c *match.Card) {

	c.Name = "Brawler Zyler"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = family.Human
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.PowerAttacker2000)

}
