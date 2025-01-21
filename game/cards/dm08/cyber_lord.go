package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// EmperorQuazla ...
func EmperorQuazla(c *match.Card) {
	c.Name = "Emperor Quazla"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	c.Use(fx.Creature, fx.Evolution, fx.Blocker, fx.When(fx.OpponentPlayedShieldTrigger, fx.DrawUpTo2))
}
