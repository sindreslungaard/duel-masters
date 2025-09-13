package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// EmperorQuazla ...
func EmperorQuazla(c *match.Card) {

	c.Name = "Emperor Quazla"
	c.Power = 5000
	c.Civ = civ.Water
	c.Family = []string{family.CyberLord}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Water}

	oppShieldTriggerCast := false

	c.Use(fx.Creature, fx.Evolution, fx.Blocker(),
		func(card *match.Card, ctx *match.Context) {
			if _, ok := ctx.Event.(*match.UntapStep); ok {
				oppShieldTriggerCast = false
			}
		},
		fx.When(fx.OppShieldTriggerCast, func(card *match.Card, ctx *match.Context) { oppShieldTriggerCast = true }),
		fx.WhenAll([]func(*match.Card, *match.Context) bool{fx.OpponentPlayedShieldTrigger, func(card *match.Card, ctx *match.Context) bool { return oppShieldTriggerCast }}, func(card *match.Card, ctx *match.Context) {
			oppShieldTriggerCast = false
			fx.DrawUpTo2(card, ctx)
		}),
	)

}
