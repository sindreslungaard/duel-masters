package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func CursedPincher(c *match.Card) {

	c.Name = "Cursed Pincher"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Blocker, fx.Slayer, fx.CantAttackPlayers, fx.CantAttackCreatures)

}
