package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CyclolinkSpectralKnight ...
func CyclolinkSpectralKnight(c *match.Card) {

	c.Name = "Cyclolink, Spectral Knight"
	c.Power = 3000
	c.Civ = civ.Light
	c.Family = []string{family.RainbowPhantom}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.WheneverThisAttacksPlayerAndIsntBlocked, fx.SearchDeckTake1Spell))

}
