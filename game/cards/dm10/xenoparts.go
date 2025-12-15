package dm10

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

// MykeesPliers ...
func MykeesPliers(c *match.Card) {

	c.Name = "Mykee's Pliers"
	c.Power = 2000
	c.Civ = civ.Fire
	c.Family = []string{family.Xenoparts}
	c.ManaCost = 4
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
		ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
			if card.Zone != match.BATTLEZONE {
				fx.FindFilter(
					card.Player,
					match.BATTLEZONE,
					func(x *match.Card) bool {
						return x.Civ == civ.Darkness || x.Civ == civ.Nature
					},
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
					return x.Civ == civ.Darkness || x.Civ == civ.Nature
				},
			).Map(func(x *match.Card) {
				x.AddUniqueSourceCondition(cnd.SpeedAttacker, true, card.ID)
			})
		})
	}))

}
