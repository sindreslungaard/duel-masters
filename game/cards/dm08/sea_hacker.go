package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// Vikorakys ...
func Vikorakys(c *match.Card) {

	c.Name = "Vikorakys"
	c.Power = 1000
	c.Civ = civ.Water
	c.Family = []string{family.SeaHacker}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Water}

	turboRush := false

	c.Use(
		fx.Creature,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) { turboRush = false }),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return turboRush }, fx.WheneverThisAttacks(fx.SearchDeckTakeXCards(1))),
	)

}
