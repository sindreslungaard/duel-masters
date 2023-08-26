package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func Gigagriff(c *match.Card) {

	c.Name = "Gigagriff"
	c.Power = 4000
	c.Civ = civ.Darkness
	c.Family = []string{family.Chimera}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker, fx.Slayer, fx.CantAttackPlayers, fx.CantAttackCreatures)
}
