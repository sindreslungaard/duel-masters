package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SalivaWorm ...
func SalivaWorm(c *match.Card) {

	c.Name = "Saliva Worm"
	c.Power = 2000
	c.Civs = []string{civ.Darkness}
	c.Family = []string{family.ParasiteWorm}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.PowerModifier = fx.WaveStrikerPower(c, 4000)

	c.Use(fx.Creature, fx.WaveStriker, fx.WaveStrikerGrant(cnd.Stealth, civ.Darkness))

}
