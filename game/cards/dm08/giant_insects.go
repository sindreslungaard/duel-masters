package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// CarboniteScarab ...
func CarboniteScarab(c *match.Card) {

	c.Name = "Carbonite Scarab"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.GiantInsect}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	turboRush := false

	attackingOpponent := false

	c.Use(
		fx.Creature,
		fx.When(fx.TurboRushCondition, func(card *match.Card, ctx *match.Context) { turboRush = true }),
		func(card *match.Card, ctx *match.Context) {
			// NOTE: @TODO currently bugged as attack can be cancelled
			// need to refactor Battle so that AttackConfirmed can tell if a creature was blocked while attacking a player
			if event, ok := ctx.Event.(*match.AttackPlayer); ok && event.CardID == c.ID {
				attackingOpponent = true
			}
		},
		fx.When(fx.EndOfMyTurn, func(card *match.Card, ctx *match.Context) {
			turboRush = false
		}),
		fx.When(func(card *match.Card, ctx *match.Context) bool { return turboRush }, func(card *match.Card, ctx *match.Context) {
			if event, ok := ctx.Event.(*match.Battle); ok {

				if event.Attacker == c && event.Blocked && attackingOpponent {
					fx.DestoryOpShield(card, ctx)
					attackingOpponent = false
				}

			}

		}))

}
