package dm07

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func ScalpelSpider(c *match.Card) {

	c.Name = "Scalpel Spider"
	c.Power = 2000
	c.Civ = civ.Darkness
	c.Family = []string{family.BrainJacker}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Darkness}

	c.Use(fx.Creature, fx.When(fx.Attacked, func(c *match.Card, ctx *match.Context) {
		c.AddCondition(cnd.Slayer, nil, c.ID)
	}))

}
