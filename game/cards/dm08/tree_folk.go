package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// SeniaOrchardAvenger ...
func SeniaOrchardAvenger(c *match.Card) {

	c.Name = "Senia, Orchard Avenger"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.TreeFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	turboRush := false

	c.Use(
		fx.Creature,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) { turboRush = false }),
		fx.When(func(c *match.Card, ctx *match.Context) bool { return turboRush }, func(card *match.Card, ctx *match.Context) {
			c.AddUniqueSourceCondition(cnd.PowerAmplifier, 5000, card.ID)
			c.AddUniqueSourceCondition(cnd.DoubleBreaker, nil, card.ID)
		}),
	)

}
