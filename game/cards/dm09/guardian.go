package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// GlenaVueleTheHypnotic ...
func GlenaVueleTheHypnotic(c *match.Card) {

	c.Name = "Glena Vuele, the Hypnotic"
	c.Power = 8500
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 5
	c.ManaRequirement = []string{civ.Light}

	oppShieldTriggerCast := false

	c.Use(fx.Creature, fx.Evolution, fx.Doublebreaker,
		func(card *match.Card, ctx *match.Context) {
			if _, ok := ctx.Event.(*match.UntapStep); ok {
				oppShieldTriggerCast = false
			}
		},
		fx.When(fx.OppShieldTriggerCast, func(card *match.Card, ctx *match.Context) { oppShieldTriggerCast = true }),
		fx.WhenAll([]func(*match.Card, *match.Context) bool{fx.OpponentPlayedShieldTrigger, func(card *match.Card, ctx *match.Context) bool { return oppShieldTriggerCast }}, func(card *match.Card, ctx *match.Context) {
			ctx.ScheduleAfter(func() { //@TODO try without ScheduleAfter
				//@TODO see binary question interaction with shield trigger pop-up
				// consider adding m.Wait(opponent) without an if check for current player turn
				oppShieldTriggerCast = false

				if fx.BinaryQuestion(
					card.Player,
					ctx.Match,
					fmt.Sprintf("%s's effect: do you want to add the top card of your deck to your shields?", card.Name)) {
					fx.TopCardToShield(card, ctx)
				}
			})
		}),
	)

}

// JilWarkaTimeGuardian ...
func JilWarkaTimeGuardian(c *match.Card) {

	c.Name = "Jil Warka, Time Guardian"
	c.Power = 2000
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.Blocker(), fx.CantAttackPlayers,
		fx.When(fx.Destroyed, func(card *match.Card, ctx *match.Context) {
			fx.Select(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: Choose up to 2 of your opponent's creatures in the battlezone and tap them.", card.Name),
				1,
				2,
				true,
			).Map(func(x *match.Card) {
				x.Tapped = true
				ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("%s was tapped by %s's effect.", x.Name, card.Name))
			})

			ctx.Match.BroadcastState()
		}))

}
