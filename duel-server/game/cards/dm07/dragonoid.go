package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func SkyCrushertheAgitator(c *match.Card) {

	c.Name = "Sky Crusher, the Agitator"
	c.Power = 4000
	c.Civ = civ.Fire
	c.Family = []string{family.Dragonoid}
	c.ManaCost = 7
	c.ManaRequirement = []string{civ.Fire}
	c.TapAbility = fx.EachPlayerDestroys1Mana

	c.Use(fx.Creature, fx.TapAbility)
}
