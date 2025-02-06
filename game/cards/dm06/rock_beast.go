package dm06

import (
	"duel-masters/game/civ"
	"duel-masters/game/cnd"
	"duel-masters/game/family"
	"duel-masters/game/fx"
	"duel-masters/game/match"
)

func RumblesaurQ(c *match.Card) {
	c.Name = "Rumblesaur Q"
	c.Power = 3000
	c.Civ = civ.Fire
	c.Family = []string{family.RockBeast, family.Survivor}
	c.ManaCost = 6
	c.ManaRequirement = []string{civ.Fire}

	c.Use(fx.Creature, fx.Survivor,
		fx.When(fx.InTheBattlezone, func(card *match.Card, ctx *match.Context) {
			ctx.Match.ApplyPersistentEffect(func(ctx2 *match.Context, exit func()) {
				if card.Zone != match.BATTLEZONE {
					fx.Find(
						card.Player,
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
					func(creature *match.Card) bool { return creature.HasFamily(family.Survivor) },
				).Map(func(x *match.Card) { x.AddUniqueSourceCondition(cnd.SpeedAttacker, true, card.ID) })
			})
		}),
	)
}
