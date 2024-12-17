package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// IllusionFish ...
func IllusionFish(c *match.Card) {

	c.Name = "Illusion Fish"
	c.Power = 3000
	c.Civ = civ.Water
	c.Family = []string{family.GelFish}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Water}

	turboRush := false

	c.Use(
		fx.Creature,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) { turboRush = false }),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return turboRush }, func(card *match.Card, ctx *match.Context) {
			c.AddUniqueSourceCondition(cnd.CantBeBlocked, nil, card.ID)
		}),
	)

}
