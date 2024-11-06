package promo

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TwisterFish ...
func TwisterFish(c *match.Card) {

	c.Name = "Twister Fish"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.ShieldTrigger)
}
