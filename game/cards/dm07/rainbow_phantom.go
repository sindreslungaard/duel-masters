package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func GeoshineSpectralKnight(c *match.Card) {

	c.Name = "Geoshine, Spectral Knight"
	c.Power = 4000
	c.Civ = civ.Light
	c.Family = []string{family.RainbowPhantom}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.WheneverThisAttacksMayTapDorFCreature())

}
