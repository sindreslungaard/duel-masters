package dm08

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// TottoPipicchi ...
func TottoPipicchi(c *match.Card) {

	c.Name = "Totto Pipicchi"
	c.Power = 1000
	c.Civ = civ.Fire
	c.Family = []string{family.FireBird}
	c.ManaCost = 3
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if card.Zone != match.BATTLEZONE {
					fx.Find(
						card.Player,
						match.BATTLEZONE,
					).Map(func(x *match.Card) {
						x.RemoveConditionBySource(card.ID)
					})

					fx.Find(
						ctx2.Match.Opponent(card.Player),
						match.BATTLEZONE,
					).Map(func(x *match.Card) {
						x.RemoveConditionBySource(card.ID)
					})

					exit()
					return
				}

				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.SharesAFamily(family.Dragons)
					},
				).Map(func(x *match.Card) {
					x.AddUniqueSourceCondition(cnd.SpeedAttacker, true, card.ID)
				})

				fx.FindFilter(
					ctx2.Match.Opponent(card.Player),
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.SharesAFamily(family.Dragons)
					},
				).Map(func(x *match.Card) {
					x.AddUniqueSourceCondition(cnd.SpeedAttacker, true, card.ID)
				})
			})
		}),
	)

}
