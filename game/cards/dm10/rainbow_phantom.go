package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// FerrosaturnSpectralKnight ...
func FerrosaturnSpectralKnight(c *match.Card) {

	c.Name = "Ferrosaturn, Spectral Knight"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.RainbowPhantom}
	c.ManaCost = 1
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers)

}
