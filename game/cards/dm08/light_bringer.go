package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// NarielTheOracle ...
func NarielTheOracle(c *match.Card) {

	c.Name = "Nariel, the Oracle"
	c.Power = 1000
	c.Civ = civ.Light
	c.Family = []string{family.LightBringer}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Light}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				fx.Find(
					card.Player,
					match.BATTLEZONE,
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				fx.Find(
					ctx.Match.Opponent(card.Player),
					match.BATTLEZONE,
				).Map(func(x *match.Card) {
					x.RemoveConditionBySource(card.ID)
				})

				exit()
				return
			}

			if _, ok := ctx2.Event.(*match.GetPowerEvent); ok {
				// to prevent infinite loop due to calling Match.GetPower() below
				return
			}

			fx.Find(
				card.Player,
				match.BATTLEZONE,
			).Map(func(x *match.Card) {
				x.RemoveConditionBySource(card.ID)
			})

			fx.Find(
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
			).Map(func(x *match.Card) {
				x.RemoveConditionBySource(card.ID)
			})

			fx.FindFilter(
				card.Player,
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return ctx2.Match.GetPower(x, false) >= 3000
				},
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.CantAttackCreatures, true, card.ID)
				x.AddUniqueSourceCondition(cnd.CantAttackPlayers, true, card.ID)
			})

			fx.FindFilter(
				ctx.Match.Opponent(card.Player),
				match.BATTLEZONE,
				func(x *match.Card) bool {
					return ctx2.Match.GetPower(x, false) >= 3000
				},
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.CantAttackCreatures, true, card.ID)
				x.AddUniqueSourceCondition(cnd.CantAttackPlayers, true, card.ID)
			})
		})
	}))

}
