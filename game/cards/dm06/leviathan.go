package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func KingTriumphant(c *match.Card) {

	c.Name = "King Triumphant"
	c.Power = 7000
	c.Civ = civ.Water
	c.Family = []string{family.Leviathan}
	c.ManaCost = 8
	c.ManaRequirement = []string{civ.Water}

	isBlocker := false

	c.Use(fx.Creature, fx.Doublebreaker,
		fx.When(ReceiveBlockerWhenOpponentPlaysCreatureOrSpell,
			func(c *match.Card, ctx *match.Context) { isBlocker = true }),
		fx.When(fx.EndOfTurn,
			func(c *match.Card, ctx *match.Context) {
				isBlocker = false
				c.RemoveConditionBySource(c.ID)
			}),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return isBlocker },
			func(c *match.Card, ctx *match.Context) {
				fx.ForceBlocker(c, ctx, c.ID)
			}),
	)
}
