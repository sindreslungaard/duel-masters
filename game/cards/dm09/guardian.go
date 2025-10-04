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
		fx.When(fx.OppShieldTriggerCast, func(card *match.Card, ctx *match.Context) {
			oppShieldTriggerCast = true
		}),
		fx.When(fx.OpponentPlayedShieldTrigger, func(card *match.Card, ctx *match.Context) {
			if oppShieldTriggerCast {
				oppShieldTriggerCast = false

				if fx.BinaryQuestion(
					card.Player,
					ctx.Match,
					fmt.Sprintf("%s's effect: do you want to add the top card of your deck to your shields?", card.Name)) {
					fx.TopCardToShield(card, ctx)
				}
			}
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
				fmt.Sprintf("%s's effect: You may choose up to 2 of your opponent's creatures in the battlezone and tap them.", card.Name),
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

// TraRionPenumbraGuardian ...
func TraRionPenumbraGuardian(c *match.Card) {

	c.Name = "Tra Rion, Penumbra Guardian"
	c.Power = 5500
	c.Civ = civ.Light
	c.Family = []string{family.Guardian}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Light}
	c.TapAbility = traRionPenumbraGuardianTapAbility

	c.Use(fx.Creature, fx.TapAbility)

}

func traRionPenumbraGuardianTapAbility(card *match.Card, ctx *match.Context) {
	family := fx.ChooseAFamily(
		card,
		ctx,
		fmt.Sprintf("%s's effect: Choose a race. At the end of this turn, untap all creatures of that race in the battlezone.", card.Name),
	)

	if family != "" {
		ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("At the end of this turn, untap all '%s' creatures in the battlezone.", family))

		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				ctx2.ScheduleAfter(func() {
					affectedCreatures := fx.FindFilter(
						card.Player,
						match.BATTLEZONE,
						func(x *match.Card) bool {
							return x.HasFamily(family)
						},
					)

					append(affectedCreatures, fx.FindFilter(
						ctx2.Match.Opponent(card.Player),
						match.BATTLEZONE,
						func(x *match.Card) bool {
							return x.HasFamily(family)
						},
					)...).Map(func(x *match.Card) {
						x.Tapped = false
						ctx2.Match.BroadcastState()
					})

					exit()
				})
			}
		})
	}
}
