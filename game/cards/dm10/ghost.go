package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// ZeroNemesisShadowOfPanic ...
func ZeroNemesisShadowOfPanic(c *match.Card) {

	c.Name = "Zero Nemesis, Shadow of Panic"
	c.Civ = civ.Darkness
	c.Power = 6000
	c.Family = []string{family.Ghost}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker, fx.WheneverOneOfMyCreaturesAttacksOppDiscardsRandom())

}
