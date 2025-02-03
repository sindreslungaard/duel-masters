package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MigaloVizierOfSpycraft ...
func MigaloVizierOfSpycraft(c *match.Card) {

	c.Name = "Migalo, Vizier of Spycraft"
	c.Power = 1500
	c.Civ = civ.Light
	c.Family = []string{family.Initiate}
	c.ManaCost = 2
	c.ManaRequirement = []string{civ.Light}

	turboRush := false

	c.Use(
		fx.Creature,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) { turboRush = false }),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return turboRush }, fx.WheneverThisAttacks(fx.ShowXShields(2))),
	)

}
