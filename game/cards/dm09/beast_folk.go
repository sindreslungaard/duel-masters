package dm09

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
	"fmt"
)

// SilvermoonTrailblazer ...
func SilvermoonTrailblazer(c *match.Card) {

	c.Name = "Silvermoon Trailblazer"
	c.Power = 3000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}
	c.TapAbility = silvermoonTrailblazerTapAbility

	c.Use(fx.Creature, fx.TapAbility)

}

func silvermoonTrailblazerTapAbility(card *match.Card, ctx *match.Context) {
	family := fx.ChooseAFamily(card, ctx, fmt.Sprintf("%s's effect: Choose a race. Creatures of that race can't be blocked by creatures that have power 3000 or less this turn.", card.Name))

	ctx.Match.ReportActionInChat(card.Player, fmt.Sprintf("Creatures of '%s' can't be blocked by creatures that have power 3000 or less this turn.", family))

	fx.FindFilter(
		card.Player,
		match.BATTLEZONE,
		func(x *match.Card) bool {
			return x.HasFamily(family)
		},
	).Map(func(x *match.Card) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
				exit()
				return
			}

			fx.CantBeBlockedByPowerUpTo3000(x, ctx2)
		})
	})
}

// StormWranglerTheFurious ...
func StormWranglerTheFurious(c *match.Card) {

	c.Name = "Storm Wrangler, the Furious"
	c.Power = 5000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.Evolution,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if card.Zone != match.BATTLEZONE {
					card.RemoveConditionBySource(card.ID)
					exit()
					return
				}

				if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
					card.RemoveConditionBySource(card.ID)
					exit()
					return
				}

				if event, ok := ctx2.Event.(*match.Battle); ok && event.Attacker == card && event.Blocked {
					event.AttackerPower += 3000                                      // for the current battle
					card.AddUniqueSourceCondition(cnd.PowerAmplifier, 3000, card.ID) // for the next potential battles in this turn
				}
			})
		}),
		fx.When(fx.AttackConfirmed, func(card *match.Card, ctx *match.Context) {
			fx.SelectFilter(
				card.Player,
				ctx.Match,
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				fmt.Sprintf("%s's effect: You may choose one of your opponent's untapped creatures that has 'Blocker'. This turn, that creature blocks this creature if able and this creature can't be blocked by other creatures.", card.Name),
				1,
				1,
				true,
				func(x *match.Card) bool {
					return !x.Tapped && x.HasCondition(cnd.Blocker)
				},
				false,
			).Map(func(x *match.Card) {
				ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
					if card.Zone != match.BATTLEZONE {
						exit()
						return
					}

					if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
						exit()
						return
					}

					fx.CantBeBlockedByOtherCreaturesBesidesX(card, ctx2, x)
				})

				ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
					if x.Zone != match.BATTLEZONE || card.Zone != match.BATTLEZONE {
						exit()
						return
					}

					if _, ok := ctx2.Event.(*match.EndOfTurnStep); ok {
						exit()
						return
					}

					if event, ok := ctx2.Event.(*match.Block); ok {
						if event.Attacker == card {
							for _, blocker := range event.Blockers {
								if blocker == x {
									// We force the opponent to block with this, i.e.
									// We cancel the Block event normal behaviour
									ctx2.InterruptFlow()

									// And we manually trigger the battle event
									// Between this creature and the selected opp blocker
									ctx2.Match.Battle(card, x, true, event.AttackedCardID == "")
								}
							}
						}
					}
				})
			})
		}))
}

// CavernRaider ...
func CavernRaider(c *match.Card) {

	c.Name = "Cavern Raider"
	c.Power = 2000
	c.Civ = civ.Nature
	c.Family = []string{family.BeastFolk}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Nature}

	c.Use(fx.Creature, fx.When(fx.WheneverThisAttacksPlayerAndIsntBlocked, fx.SearchDeckTake1Creature))

}
