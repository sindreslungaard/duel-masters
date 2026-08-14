package dm11

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SapianTarkFlameDervish ...
func SapianTarkFlameDervish(c *match.Card) {

	c.Name = "Sapian Tark, Flame Dervish"
	c.Power = 2000
	c.Civs = []string{civ.Fire}
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.PowerModifier = fx.WaveStrikerPower(c, 4000)

	c.Use(fx.Creature, fx.WaveStriker, fx.WaveStrikerGrant(cnd.AttackUntapped, true))

}
